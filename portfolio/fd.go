package portfolio

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/Samar2170/portfolio-manager/account"
	"github.com/Samar2170/portfolio-manager/securities"
)

type FDFile struct {
	Id       uint
	User     account.User `gorm:"foreignKey:UserId"`
	UserId   uint
	FileName string
	FilePath string
	Parsed   bool
}

type FDHolding struct {
	ID             uint                    `gorm:"primaryKey"`
	FixedDeposit   securities.FixedDeposit `gorm:"foreignKey:FixedDepositId"`
	FixedDepositId uint
	BankAccount    account.BankAccount `gorm:"foreignKey:BankAccountId"`
	BankAccountId  uint
}

func (fdh *FDHolding) create() error {
	err := db.Create(&fdh).Error
	return err
}

func (fdf FDFile) Create() error {
	err := db.Create(&fdf).Error
	return err
}
func getFDFileById(fileId uint) (FDFile, error) {
	var fdf FDFile
	err := db.First(&fdf, "id = ?", fileId).Error
	return fdf, err
}

func CreateFDHolding(bankName string, amount, mtAmount, ipRate float64, ipfreq string, startDate, ipDate, mtDate, accNumber string) (FDHolding, error) {
	bankNameCap := strings.ToUpper(bankName)
	if _, ok := BankNames[bankNameCap]; !ok {
		return FDHolding{}, errors.New("unknown bank name")
	}
	if ipRate < 0 || ipRate > 0.99 {
		return FDHolding{}, errors.New("ip rate not in valid range")
	}
	if _, ok := ValidIPFreqs[ipfreq]; !ok {
		return FDHolding{}, errors.New("ip freq not valid, should be A, M, MT, Q, SA")
	}

	StartDateParsed, err := time.Parse(DtFormat, startDate)
	if err != nil {
		return FDHolding{}, errors.New("startDate date not valid")
	}
	IPDateParsed, err := time.Parse(DtFormat, ipDate)
	log.Println(IPDateParsed, err)
	if err != nil {
		return FDHolding{}, errors.New("ip date not valid")
	}
	MtDateParsed, err := time.Parse(DtFormat, mtDate)
	if err != nil {
		return FDHolding{}, errors.New("mtDate date not valid")
	}
	fd := securities.FixedDeposit{
		// BankName:  bankNameCap,
		Amount:    amount,
		IPRate:    ipRate,
		IPFreq:    ipfreq,
		IPDate:    IPDateParsed,
		StartDate: StartDateParsed,
		MtDate:    MtDateParsed,
		MtAmount:  mtAmount,
	}
	err = fd.Create()
	if err != nil {
		log.Println(err)
		return FDHolding{}, err
	}
	ba, err := account.GetBankAccountByNumber(accNumber)
	if err != nil {
		return FDHolding{}, errors.New("bank account number not valid")
	}
	fdh := FDHolding{
		FixedDepositId: fd.ID,
		BankAccount:    ba,
	}
	err = fdh.create()
	if err != nil {
		return FDHolding{}, err
	}
	return fdh, nil

}

// func ParseFDFile(fileId uint) error {
// 	fdFile, err := getFDFileById(fileId)
// 	if err != nil {
// 		return err
// 	}
// }
