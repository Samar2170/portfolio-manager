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
