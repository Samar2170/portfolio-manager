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

type UnlistedNCDHolding struct {
	*gorm.Model
	Id               uint
	UnlistedNCD      securities.UnlistedNCD `gorm:"foreignKey:UnlistedNCDId"`
	UnlistedNCDId    uint
	Quantity         float64
	Price            float64
	TradeDate        time.Time
	GeneralAccount   account.GeneralAccount `gorm:"foreignKey:GeneralAccountId"`
	GeneralAccountId uint
}

func (ulncdh UnlistedNCDHolding) Create() error {
	err := db.Create(&ulncdh).Error
	return err
}
func CreateUnlistedNCDHolding(symbol, quantity, price, tradeDate string, userId uint) (UnlistedNCDHolding, error) {
	var tradeDateParsed time.Time
	unlisteNcd, err := securities.GetUnlistedNCDBySymbol(symbol)
	if err != nil {
		return UnlistedNCDHolding{}, errors.New("unlisted ncd not found. Register first using register-listed-ncd")
	}
	quantityFloat, err := strconv.ParseFloat(quantity, 64)
	if err != nil {
		return UnlistedNCDHolding{}, errors.New("quantity should be a number")
	}
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return UnlistedNCDHolding{}, errors.New("price should be a number")
	}
	if tradeDate == "" {
		tradeDateParsed = time.Now()
	} else {
		tradeDateParsed, err = time.Parse(DtFormat, tradeDate)
		if err != nil {
			return UnlistedNCDHolding{}, errors.New("trade date should be in format 2022-11-22")
		}
	}
	ga, err := account.GetGeneralAccountByUser(userId)
	if err != nil {
		return UnlistedNCDHolding{}, fmt.Errorf("user Id %d doesnt have a general account", userId)
	}
	ulncdh := UnlistedNCDHolding{
		UnlistedNCD:    unlisteNcd,
		Quantity:       quantityFloat,
		Price:          priceFloat,
		TradeDate:      tradeDateParsed,
		GeneralAccount: ga,
	}
	err = ulncdh.Create()
	if err != nil {
		return UnlistedNCDHolding{}, fmt.Errorf("something went wrong")
	}
	return ulncdh, nil
}
