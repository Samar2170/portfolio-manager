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

func CalculateNextIPDate(fd FixedDeposit) time.Time {
	var nextIPDate time.Time
	today := time.Now()
	ipdate := fd.IPDate
	var currIpdate time.Time
	if ipdate.Before(today) {
		switch fd.IPFreq {
		case "A":
			currIpdate = time.Date(today.Year(), ipdate.Month(), ipdate.Day(), 0, 0, 0, 0, time.Local)
			if currIpdate.Before(today) {
				nextIPDate = currIpdate.AddDate(1, 0, 0)
			} else {
				nextIPDate = currIpdate
			}

		case "M":
			currIpdate = time.Date(today.Year(), today.Month(), ipdate.Day(), 0, 0, 0, 0, time.Local)
			nextIPDate = currIpdate.AddDate(0, 1, 0)
		case "Q":
			nextIPDate = utils.GetNextQuarter(today)
		case "MT":
			nextIPDate = fd.MtDate
		case "SA":
			nextIPDate = utils.GetNextHY(today)
			// case "QAD":
			// 	nextIPDate = utils.GetNextQuarter(today)
			// case "SAD":
			// 	nextIPDate = utils.GetNextHY(today)
		}
	} else {
		nextIPDate = currIpdate
	}
	return nextIPDate
}

func UpdateNextIPDatesFDs() error {
	var fds []FixedDeposit

	err := db.Where("next_ip_date < CURRENT_DATE").Find(&fds).Error
	if err != nil {
		return err
	}

	for _, fd := range fds {
		nipd := CalculateNextIPDate(fd)
		err = db.Model(&fd).Where("id = ?", fd.ID).Update("next_ip_date", nipd).Error
	}
	return err
}

func (fd *FixedDeposit) Create() error {
	return db.Create(&fd).Error
}
