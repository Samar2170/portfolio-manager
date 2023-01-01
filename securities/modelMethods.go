package securities

import "fmt"

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

func (mf MutualFund) create() error {
	err := db.Create(&mf).Error
	return err
}

func GetAllMutualFunds() ([]MutualFund, error) {
	var mfs []MutualFund
	err := db.Find(&mfs).Error
	return mfs, err
}

func GetStockBySymbol(symbol string) (Stock, error) {
	var stock Stock
	err := db.First(&stock, "symbol = ?", symbol).Error
	return stock, err
}
