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

type MFTrade struct {
	*gorm.Model
	TradeDate      time.Time
	MF             securities.MutualFund `gorm:"foreignKey:MFId"`
	MFId           uint
	Quantity       float64
	TradeType      string
	Price          float64
	DematAccount   account.DematAccount `gorm:"foreignKey:DematAccountId"`
	DematAccountId uint
}
type MFHolding struct {
	*gorm.Model
	MF             securities.MutualFund `gorm:"foreignKey:MFId"`
	MFId           uint
	Quantity       float64
	Price          float64
	DematAccount   account.DematAccount `gorm:"foreignKey:DematAccountId"`
	DematAccountId uint
}

type MFFile struct {
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

func (mft MFTrade) Create() error {
	err := db.Create(&mft).Error
	return err
}
func (mfh MFHolding) Create() error {
	err := db.Create(&mfh).Error
	return err
}

func NewMFTrade(mfId uint, tradeType, dematAccCode string, quantity, price float64, tradeDate time.Time) (*MFTrade, error) {
	tradeTypeCap := strings.ToUpper(tradeType)
	if tradeTypeCap != "BUY" && tradeTypeCap != "SELL" {
		return &MFTrade{}, errors.New("trade type should be buy or sell")
	}
	mf, err := securities.GetMutualFundById(mfId)
	if err != nil {
		return &MFTrade{}, errors.New("unable to find mutual fund please check the symbol provided")
	}
	da, err := account.GetDematAccountByCode(dematAccCode)
	if err != nil {
		return &MFTrade{}, errors.New("demat Account Code is not valid.Please Check")
	}
	return &MFTrade{
		MF:           mf,
		Quantity:     quantity,
		Price:        price,
		DematAccount: da,
		TradeType:    tradeTypeCap,
	}, nil
}

func GetMFHolding(mf securities.MutualFund, dematAccount account.DematAccount) (MFHolding, error) {
	var mfh MFHolding
	err := db.Where("mf_id = ? AND demat_account_id = ?", mf.ID, dematAccount.ID).First(&mfh).Error
	return mfh, err

}

func checkMFHoldings(mf securities.MutualFund, dematAccount account.DematAccount) bool {
	var count int64
	db.Model(&MFHolding{}).Where("mf_id = ?", mf.ID).Where("demat_account_id = ?", dematAccount.ID).Count(&count)
	return count > 0
}

func RegisterMFTrade(mft MFTrade) error {
	holdingExists := checkMFHoldings(mft.MF, mft.DematAccount)
	if holdingExists {
		mfh, err := GetMFHolding(mft.MF, mft.DematAccount)
		if err != nil {
			return err
		}
		if mft.TradeType == "BUY" {

			newPrice := ((mfh.Quantity * mfh.Price) + (mft.Price * mft.Quantity)) / (float64(mft.Quantity) + float64(mfh.Quantity))
			newQuantity := mfh.Quantity + mft.Quantity
			mft.Create()
			db.Model(&mfh).Where("mf_id = ? ", mft.MF.ID).Where("demat_account_id = ?", mft.DematAccount.ID).Updates(map[string]interface{}{
				"quantity": newQuantity, "price": newPrice})
		}
		if mft.TradeType == "SELL" {
			fmt.Println(mft.Quantity, mfh)
			if mft.Quantity > mfh.Quantity {
				return errors.New("cant sell more than you own")
			}
			mft.Create()
			mfh.Quantity -= mft.Quantity
			db.Model(&mfh).Where("mf_id = ? ", mft.MF.ID).Where("demat_account_id = ?", mft.DematAccount.ID).Update("quantity", mfh.Quantity)
		}
	} else {
		if mft.TradeType == "BUY" {
			mfh := MFHolding{
				MF:           mft.MF,
				DematAccount: mft.DematAccount,
				Quantity:     mft.Quantity,
				Price:        mft.Price,
			}
			err := mft.Create()
			if err != nil {
				return err
			}

			err = mfh.Create()
			if err != nil {
				return err
			}

		} else {
			return errors.New("no Holding found to sell")
		}
	}
	return nil
}

func (mff MFFile) Create() error {
	err := db.Create(&mff).Error
	return err
}

func getMFFIleById(fileId uint) (MFFile, error) {
	var mff MFFile
	err := db.First(&mff, "id = ?", fileId).Error
	return mff, err
}

func CreateMFTrade(mfId, dematAccCode, quantity, price, tradeType, tradeDate string) (MFTrade, error) {
	var err error
	var tradeDateParsed time.Time
	mfIdInt, err := strconv.ParseInt(mfId, 10, 64)
	if err != nil {
		return MFTrade{}, errors.New("mutual_fund_id should be a number")
	}

	if tradeDate == "" {
		tradeDateParsed = time.Now()
	} else {
		tradeDateParsed, err = time.Parse(DtFormat, tradeDate)
		if err != nil {
			return MFTrade{}, errors.New("trade date should be in format 2022-11-22")
		}
	}
	quantityFloat, err := strconv.ParseFloat(quantity, 64)
	if err != nil {
		return MFTrade{}, errors.New("quantity should be a number")
	}
	priceFloat, err := strconv.Atoi(price)
	if err != nil {
		return MFTrade{}, errors.New("price should be a number")
	}
	mfTrade, err := NewMFTrade(uint(mfIdInt), tradeType, dematAccCode, quantityFloat, float64(priceFloat), tradeDateParsed)
	if err != nil {
		return MFTrade{}, err
	}
	err = RegisterMFTrade(*mfTrade)
	if err != nil {
		return MFTrade{}, err
	}
	return *mfTrade, nil
}

func ParseMFFIle(fileId uint) error {
	fileData, err := getMFFIleById(fileId)
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
		_, err := CreateMFTrade(record[0], record[1], record[2], record[3], record[4], record[5])
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
