package account

import (
	"errors"
	"strings"
	"time"

	"github.com/Samar2170/portfolio-manager/securities"
	"gorm.io/gorm"
)

type User struct {
	*gorm.Model
	Username string `gorm:"type:varchar(100);uniqueIndex"`
	Password string `gorm:"type:varchar(100)"`
	Email    string `gorm:"type:varchar(100);uniqueIndex"`
}

type GeneralAccount struct {
	*gorm.Model
	Code   string `gorm:"unique_index"`
	User   User   `gorm:"foreignKey:UserId"`
	UserId uint
}

type DematAccount struct {
	*gorm.Model
	Code   string `gorm:"uniqueIndex"`
	User   User   `gorm:"foreignKey:UserId"`
	UserId uint
	Broker string
}

type BankAccount struct {
	*gorm.Model
	AccountNo string `gorm:"uniqueIndex"`
	User      User   `gorm:"foreignKey:UserId"`
	UserId    uint
	Bank      string
}

type StockTrade struct {
	TradeDate    time.Time
	Stock        securities.Stock
	Quantity     uint
	Price        float64
	TradeType    string
	DematAccount DematAccount
}

type StockHolding struct {
	Stock        securities.Stock
	Quantity     uint
	Price        float64
	DematAccount DematAccount
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
	da, err := GetDematAccountByCode(dematAccountCode)

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

func (u *User) GetOrCreate() (User, error) {
	var existingUser User
	_ = db.Where("username = ?", u.Username).First(&existingUser).Error
	if existingUser.Username == u.Username {
		return User{}, errors.New("username already exists")
	}
	err := db.Create(&u).Error
	return *u, err
}

func GetUserById(userId uint) (User, error) {
	var user User
	err := db.First("id = ?", userId).Error
	return user, err
}

func GetUserByUsername(username string) (User, error) {
	var user User
	err := db.First("username = ", username).Error
	return user, err
}

func (ga *GeneralAccount) GetOrCreate() (GeneralAccount, error) {
	err := db.Create(&ga).Error
	return *ga, err
}

func (da *DematAccount) GetOrCreate() (DematAccount, error) {
	err := db.Create(&da).Error
	return *da, err
}
func (ba *BankAccount) GetOrCreate() (BankAccount, error) {
	err := db.Create(&ba).Error
	return *ba, err
}

func GetDematAccountByCode(code string) (DematAccount, error) {
	var da DematAccount
	err := db.First("code = ?", code).Error
	return da, err
}
