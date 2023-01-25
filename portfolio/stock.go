package portfolio

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Samar2170/portfolio-manager/account"
	"github.com/Samar2170/portfolio-manager/securities"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type StockTrade struct {
	*gorm.Model
	TradeDate      time.Time
	Stock          securities.Stock `gorm:"foreignKey:StockId"`
	StockId        uint
	Quantity       uint
	Price          float64
	TradeType      string
	DematAccount   account.DematAccount `gorm:"foreignKey:DematAccountId"`
	DematAccountId uint
}

type StockHolding struct {
	*gorm.Model
	Stock          securities.Stock `gorm:"foreignKey:StockId"`
	StockId        uint
	Quantity       uint
	Price          float64
	DematAccount   account.DematAccount `gorm:"foreignKey:DematAccountId"`
	DematAccountId uint
}

type StockFile struct {
	*gorm.Model
	Id         uint
	User       account.User `gorm:"foreignKey:UserId"`
	UserId     uint
	FileName   string
	FilePath   string
	Parsed     bool
	RowsFailed pq.Int64Array  `gorm:"type:integer[]"`
	RowErrors  pq.StringArray `gorm:"type:varchar[]"`
}

func (st StockTrade) Create() error {
	err := db.Create(&st).Error
	return err
}

func (sh StockHolding) Create() error {
	err := db.Create(&sh).Error
	return err
}

func NewStockTrade(symbol, tradeType, dematAccountCode string, quantity uint, price float64, tradeDate time.Time) (*StockTrade, error) {
	tradeTypeCap := strings.ToUpper(tradeType)
	if tradeTypeCap != "BUY" && tradeTypeCap != "SELL" {
		return &StockTrade{}, errors.New("trade type should be buy or sell")
	}
	stock, err := securities.GetStockBySymbol(symbol)
	if err != nil {
		return &StockTrade{}, errors.New("unable to find stock please check the symbol provided")
	}
	da, err := account.GetDematAccountByCode(dematAccountCode)
	if err != nil {
		return &StockTrade{}, errors.New("demat Account Code is not valid.Please Check")
	}
	return &StockTrade{
		Stock:        stock,
		TradeDate:    tradeDate,
		Quantity:     quantity,
		Price:        price,
		TradeType:    tradeTypeCap,
		DematAccount: da,
	}, nil

}

func GetStockHolding(stock securities.Stock, dematAccount account.DematAccount) (StockHolding, error) {
	var sth StockHolding
	err := db.Where("stock_id = ? AND demat_account_id = ?", stock.ID, dematAccount.ID).First(&sth).Error
	return sth, err
}

func checkStockHoldings(stock securities.Stock, dematAcc account.DematAccount) bool {
	var count int64

	db.Model(&StockHolding{}).Where("stock_id = ? ", stock.ID).Where("demat_account_id = ?", dematAcc.ID).Count(&count)
	return count > 0
}

func RegisterTrade(nst StockTrade) error {
	holdingExists := checkStockHoldings(nst.Stock, nst.DematAccount)
	if holdingExists {
		sth, err := GetStockHolding(nst.Stock, nst.DematAccount)
		if err != nil {
			return err
		}
		if nst.TradeType == "BUY" {

			newPrice := ((float64(sth.Quantity) * sth.Price) + (nst.Price * float64(nst.Quantity))) / (float64(nst.Quantity) + float64(sth.Quantity))
			newQuantity := sth.Quantity + nst.Quantity
			nst.Create()
			db.Model(&sth).Where("stock_id = ? ", nst.Stock.ID).Where("demat_account_id = ?", nst.DematAccount.ID).Updates(map[string]interface{}{
				"quantity": newQuantity, "price": newPrice})
		}
		if nst.TradeType == "SELL" {
			fmt.Println(nst.Quantity, sth)
			if nst.Quantity > sth.Quantity {
				return errors.New("cant sell more than you own")
			}
			nst.Create()
			sth.Quantity -= nst.Quantity
			db.Model(&sth).Where("stock_id = ? ", nst.Stock.ID).Where("demat_account_id = ?", nst.DematAccount.ID).Update("quantity", sth.Quantity)
		}
	} else {
		if nst.TradeType == "BUY" {
			sh := StockHolding{
				Stock:        nst.Stock,
				DematAccount: nst.DematAccount,
				Quantity:     nst.Quantity,
				Price:        nst.Price,
			}
			err := nst.Create()
			if err != nil {
				return err
			}

			err = sh.Create()
			if err != nil {
				return err
			}

		} else {
			return errors.New("no Holding found to sell")
		}
	}
	return nil
}

func (sf StockFile) Create() error {
	err := db.Create(&sf).Error
	return err
}
func GetStockFileById(sfId uint) (StockFile, error) {
	var stockFile StockFile
	err := db.First(&stockFile, "id = ?", sfId).Error
	return stockFile, err
}

func CreateStockTrade(symbol, dematAccCode, quantity, price, tradeType, tradeDate string) (StockTrade, error) {
	var err error
	var tradeDateParsed time.Time
	symbol = strings.ToUpper(symbol)
	if err != nil {
		return StockTrade{}, err
	}
	if tradeDate == "" {
		tradeDateParsed = time.Now()
	} else {
		tradeDateParsed, err = time.Parse(DtFormat, tradeDate)
		if err != nil {
			return StockTrade{}, errors.New("trade date should be in format 2022-11-22")
		}
	}
	quantityFloat, err := strconv.ParseFloat(quantity, 64)
	if err != nil {
		return StockTrade{}, errors.New("quantity should be a number")
	}
	priceFloat, err := strconv.Atoi(price)
	if err != nil {
		return StockTrade{}, errors.New("price should be a number")
	}
	stockTrade, err := NewStockTrade(symbol, tradeType, dematAccCode, uint(quantityFloat), float64(priceFloat), tradeDateParsed)
	if err != nil {
		return StockTrade{}, errors.New(err.Error())
	}
	err = RegisterTrade(*stockTrade)
	return *stockTrade, err
}
func ParseStockFile(fileId uint) error {
	fileData, err := GetStockFileById(fileId)
	if err != nil {
		return err
	}
	file, err := os.Open(fileData.FilePath)
	if err != nil {
		return err
	}
	r := csv.NewReader(file)
	if _, err := r.Read(); err != nil {
		return err
	}
	failedRows := []int64{}
	errorRows := []string{}
	records, err := r.ReadAll()
	if err != nil {
		return err
	}
	for i, record := range records {
		_, err := CreateStockTrade(record[0], record[1], record[2], record[3], record[4], record[5])
		if err != nil {
			errorRows = append(errorRows, err.Error())
			failedRows = append(failedRows, int64(i))
		}
	}
	fileData.RowErrors = errorRows
	fileData.RowsFailed = failedRows
	fileData.Parsed = true
	db.Save(&fileData)
	return nil
}
