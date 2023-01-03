package portfolio

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Samar2170/portfolio-manager/account"
	"github.com/Samar2170/portfolio-manager/securities"
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
