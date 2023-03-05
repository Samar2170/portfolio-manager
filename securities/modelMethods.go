package securities

import (
	"fmt"

	"github.com/Samar2170/portfolio-manager/utils"
)

func (s Stock) create() error {
	fmt.Println(s)
	err := db.Create(&s).Error
	return err
}

func GetAllStocks() ([]Stock, error) {
	var stocks []Stock
	err := db.Find(&stocks).Error
	return stocks, err
}
func GetStocksSymbols() ([]string, error) {
	var stocks []Stock
	var symbols []string

	err := db.Find(&stocks).Error
	if err != nil {
		return symbols, err
	}
	for _, stock := range stocks {
		symbols = append(symbols, stock.Symbol)
	}
	return symbols, nil

}

func GetStockBySymbol(symbol string) (Stock, error) {
	var stock Stock
	err := db.First(&stock, "symbol = ?", symbol).Error
	return stock, err
}

func SearchStockSymbol(symbol string, pagination utils.Pagination) (*utils.Pagination, []*Stock, error) {
	var stocks []*Stock
	err := db.Scopes(utils.Paginate(stocks, &pagination, db)).Where("symbol ILIKE ?", "%"+symbol+"%").Find(&stocks).Error
	return &pagination, stocks, err
}
