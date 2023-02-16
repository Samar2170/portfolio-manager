package securities

import (
	"time"

	"gorm.io/gorm"
)

type Stock struct {
	*gorm.Model
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"uniqueIndex"`
	Symbol       string `gorm:"uniqueIndex"`
	SecurityCode string `gorm:"unique"`
	Exchange     string
	Industry     string
}

type Crypto struct {
	*gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"uniqueIndex"`
	Symbol   string `gorm:"uniqueIndex"`
	Exchange string
}

type FixedDeposit struct {
	*gorm.Model
	ID uint `gorm:"primaryKey"`
	// BankName   string
	Amount          float64
	IPRate          float64
	IPFreq          string
	IPDate          time.Time
	StartDate       time.Time
	MtDate          time.Time
	MtAmount        float64
	NextIPDate      time.Time
	AccruedInterest float64
}
