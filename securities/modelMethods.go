package securities

import (
	"fmt"
	"time"

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

func SearchStockSymbol(symbol string, pagination utils.Pagination) (*utils.Pagination, []*Stock, error) {
	var stocks []*Stock
	err := db.Scopes(utils.Paginate(stocks, &pagination, db)).Where("symbol ILIKE ?", "%"+symbol+"%").Find(&stocks).Error
	return &pagination, stocks, err
}

func SearchMutualFunds(symbol string, pagination utils.Pagination) (*utils.Pagination, []*MutualFund, error) {
	var mfs []*MutualFund
	err := db.Scopes(utils.Paginate(mfs, &pagination, db)).Where("scheme_nav_name ILIKE ?", "%"+symbol+"%").Find(&mfs).Error
	return &pagination, mfs, err
}

func GetMutualFundById(mfId uint) (MutualFund, error) {
	var mf MutualFund
	err := db.First(&mf, "id = ?", mfId).Error
	return mf, err
}
func GetFDByID(id uint) (FixedDeposit, error) {
	var fd FixedDeposit
	err := db.First(&fd, "id = ?", id).Error
	return fd, err
}

func (fd FixedDeposit) CalculateNextIPDate() time.Time {
	var nextIPDate time.Time
	today := time.Now()
	if today.Before(fd.IPDate) {
		switch fd.IPFreq {
		case "A":
			nextIPDate = fd.IPDate.AddDate(1, 0, 0)
		case "M":
			nextIPDate = fd.IPDate.AddDate(0, 1, 0)
		case "Q":
			nextIPDate = fd.IPDate.AddDate(0, 3, 0)
		case "MT":
			nextIPDate = fd.MtDate
		case "SA":
			nextIPDate = fd.IPDate.AddDate(0, 6, 0)
		}
	}
	return nextIPDate
}

func UpdateNextIPDatesFDs() error {
	var fds []FixedDeposit
	err := db.Where("next_ip_date < CURRENT_DATE").Find(&fds).Error
	if err != nil {
		return err
	}
	for i, fd := range fds {
		fmt.Println(i, fd)
	}
	return nil
}
