package portfolio

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Samar2170/portfolio-manager/account"
	"github.com/Samar2170/portfolio-manager/securities"
)

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
