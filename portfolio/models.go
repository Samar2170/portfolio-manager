package portfolio

import (
	"time"

	"github.com/Samar2170/portfolio-manager/account"
	"github.com/Samar2170/portfolio-manager/securities"
	"gorm.io/gorm"
)

type StockTrade struct {
	*gorm.Model
	TradeDate      time.Time
	Stock          securities.Stock `gorm:"foreignKey:StockId"`
	StockId        uint
	Quantity       uint
	Price          float64
	TradeType      string
	DematAccount   account.DematAccount `gorm:"foreignKey:DematAccountId"`
	DematAccountId uint
}

type StockHolding struct {
	*gorm.Model
	Stock          securities.Stock `gorm:"foreignKey:StockId"`
	StockId        uint
	Quantity       uint
	Price          float64
	DematAccount   account.DematAccount `gorm:"foreignKey:DematAccountId"`
	DematAccountId uint
}

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
	Quantity       uint
	Price          float64
	DematAccount   account.DematAccount `gorm:"foreignKey:DematAccountId"`
	DematAccountId uint
}
