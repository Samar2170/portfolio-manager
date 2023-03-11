package portfolio

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Samar2170/portfolio-manager/account"
	"github.com/Samar2170/portfolio-manager/securities"
	"gorm.io/gorm"
)

type ListedNCDHolding struct {
	*gorm.Model
	Id             uint
	ListedNCD      securities.ListedNCD `gorm:"foreignKey:ListedNCDId"`
	ListedNCDId    uint
	Quantity       float64
	DirtyPrice     float64
	TradeDate      time.Time
	DematAccount   account.DematAccount `gorm:"foreignKey:DematAccountId"`
	DematAccountId uint
}

func (lncdh ListedNCDHolding) Create() error {
	err := db.Create(&lncdh).Error
	return err
}
func CreateListedNCDHolding(symbol, quantity, price, tradeDate, dematAccountCode string, userId uint) (ListedNCDHolding, error) {
	var tradeDateParsed time.Time
	listeNcd, err := securities.GetListedNCDBySymbol(symbol)
	if err != nil {
		return ListedNCDHolding{}, errors.New("listed ncd not found. Register first using register-listed-ncd")
	}
	quantityFloat, err := strconv.ParseFloat(quantity, 64)
	if err != nil {
		return ListedNCDHolding{}, errors.New("quantity should be a number")
	}
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return ListedNCDHolding{}, errors.New("price should be a number")
	}
	if tradeDate == "" {
		tradeDateParsed = time.Now()
	} else {
		tradeDateParsed, err = time.Parse(DtFormat, tradeDate)
		if err != nil {
			return ListedNCDHolding{}, errors.New("trade date should be in format 2022-11-22")
		}
	}
	check := account.CheckDematAccountAndUserId(userId, dematAccountCode)
	if !check {
		return ListedNCDHolding{}, fmt.Errorf("demat %s & User Id %d dont match", dematAccountCode, userId)
	}
	dematAccount, _ := account.GetDematAccountByCode(dematAccountCode)
	lncdh := ListedNCDHolding{
		ListedNCD:    listeNcd,
		DematAccount: dematAccount,
		Quantity:     quantityFloat,
		DirtyPrice:   priceFloat,
		TradeDate:    tradeDateParsed,
	}
	err = lncdh.Create()
	if err != nil {
		return ListedNCDHolding{}, fmt.Errorf("something went wrong")
	}
	return lncdh, nil
}
func GetListedNCDByUser(userId uint) []ListedNCDHolding {
	var holdings []ListedNCDHolding
	dematIds, _ := account.GetDematAccountIdsByUser(userId)
	db.Joins("ListedNCD").Find(&holdings, "demat_account_id IN ?", dematIds)
	return holdings
}

func (lncdh ListedNCDHolding) getHoldings() HoldingSecurity {
	return HoldingSecurity{
		Name:         lncdh.ListedNCD.Name,
		CurrentValue: lncdh.DirtyPrice * lncdh.DirtyPrice,
		Invested:     lncdh.Quantity * lncdh.DirtyPrice,
		Category:     "Listed-NCD",
	}
}
