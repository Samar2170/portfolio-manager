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

type ListedNCD struct {
	*gorm.Model
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"uniqueIndex"`
	Symbol        string `gorm:"uniqueIndex"`
	SecurityCode  string `gorm:"unique"`
	Exchange      string
	IPFreq        string
	IPRate        float64
	IPDate        time.Time
	MtDate        time.Time
	FaceValue     float64
	MaturityValue float64
}

type Crypto struct {
	*gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"uniqueIndex"`
	Symbol   string `gorm:"uniqueIndex"`
	Exchange string
}

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

type UnlistedNCD struct {
	*gorm.Model
	ID            uint `gorm:"primaryKey"`
	Name          uint `gorm:"unique"`
	IPFreq        string
	IPDate        time.Time
	MtDate        time.Time
	IPRate        float64
	FaceValue     float64
	MaturityValue float64
}

type FixedDeposit struct {
	*gorm.Model
	ID         uint   `gorm:"primaryKey"`
	BankName   string `gorm:"unique"`
	Amount     float64
	IPRate     float64
	IPFreq     string
	IPDate     time.Time
	StartDate  time.Time
	MtDate     time.Time
	MtAmount   float64
	NextIPDate time.Time
}
