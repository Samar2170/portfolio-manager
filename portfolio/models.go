package portfolio

import (
	"errors"
	"strings"
	"time"

	"github.com/Samar2170/portfolio-manager/account"
	"github.com/Samar2170/portfolio-manager/securities"
)

type StockTrade struct {
	TradeDate    time.Time
	Stock        securities.Stock
	Quantity     uint
	Price        float64
	TradeType    string
	DematAccount account.DematAccount
}

type StockHolding struct {
	Stock        securities.Stock
	Quantity     uint
	Price        float64
	DematAccount account.DematAccount
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

// func checkHolding(userId uint, symbol string) bool {

// }
