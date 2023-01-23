package securities

import (
	"github.com/Samar2170/portfolio-manager/utils"
	"gorm.io/gorm"
)

type MutualFund struct {
	*gorm.Model
	ID             uint `gorm:"primaryKey"`
	AMC            string
	Code           string
	SchemeName     string
	SchemeType     string
	SchemeCategory string
	SchemeNAVName  string `gorm:"unique"`
	Category       string
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
