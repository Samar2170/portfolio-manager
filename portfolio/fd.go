package portfolio

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"net/textproto"
	"strings"
	"sync"
	"time"

	"github.com/Samar2170/portfolio-manager/account"
	"github.com/Samar2170/portfolio-manager/securities"
	"github.com/jordan-wright/email"
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

//	func ParseFDFile(fileId uint) error {
//		fdFile, err := getFDFileById(fileId)
//		if err != nil {
//			return err
//		}
//	}
type InterestDueFDResult struct {
	Bank       string
	Amount     float64
	NextIPDate time.Time
	UserId     uint
	Email      string
}

func FindInterestDueFD() error {
	t := time.Now()
	t2 := t.AddDate(0, 0, 7)
	var results []InterestDueFDResult
	db.Raw("SELECT bank_accounts.bank, bank_accounts.user_id,users.email, fixed_deposits.amount, fixed_deposits.next_ip_date FROM fd_holdings LEFT JOIN fixed_deposits ON fixed_deposits.id = fixed_deposit_id LEFT JOIN bank_accounts ON bank_accounts.id = bank_account_id LEFT JOIN users ON users.id = bank_accounts.user_id WHERE fixed_deposits.next_ip_date >= ? AND fixed_deposits.next_ip_date <= ?", t, t2).Scan(&results)

	emails := []*email.Email{}
	for _, result := range results {
		msgText := fmt.Sprintf("Interest of FD of %d at Bank %s is Due on %s", int(result.Amount), result.Bank, result.NextIPDate.Format(time.ANSIC))
		emails = append(emails, &email.Email{
			To:      []string{result.Email},
			From:    EMAILID,
			Subject: "Interest Due",
			Text:    []byte(msgText),
			HTML:    []byte(""),
			Headers: textproto.MIMEHeader{},
		})
	}
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         "localhost",
	}
	var Wg sync.WaitGroup
	Wg.Add(len(emails))

	for _, msg := range emails {
		fmt.Println(msg)
		go func(msg *email.Email) {
			defer Wg.Done()
			err := msg.SendWithStartTLS(SmtpAddressWPort, smtp.PlainAuth("", EMAILID, EMAILPASSWORD, SmtpAddress), tlsconfig)
			if err != nil {
				log.Println(err)
			}
		}(msg)
	}
	Wg.Wait()
	return nil
}
