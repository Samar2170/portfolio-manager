package portfolio

import (
	"crypto/tls"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"net/textproto"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Samar2170/portfolio-manager/account"
	"github.com/Samar2170/portfolio-manager/securities"
	"github.com/jordan-wright/email"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type FDFile struct {
	*gorm.Model
	Id         uint
	User       account.User `gorm:"foreignKey:UserId"`
	UserId     uint
	FileName   string
	FilePath   string
	Parsed     bool
	RowsFailed pq.Int64Array  `gorm:"type:integer[]"`
	RowErrors  pq.StringArray `gorm:"type:varchar[]"`
}

type FDHolding struct {
	*gorm.Model
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

func (fdf FDFile) Create() (FDFile, error) {
	err := db.Create(&fdf).Error
	return fdf, err
}
func GetFDFileById(fileId uint) (FDFile, error) {
	var fdf FDFile
	err := db.First(&fdf, "id = ?", fileId).Error
	return fdf, err
}

func CreateFDHolding(bankName string, amount, mtAmount, ipRate float64, ipfreq string, startDate, ipDate, mtDate, accNumber string, userId uint) (FDHolding, error) {
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
	check := account.CheckBankAccountAndUserId(userId, accNumber)
	if !check {
		return FDHolding{}, errors.New("bank account and user do not match")
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

func ParseFDFile(fileId uint) error {
	// work pending on this/ ? should we go with one master file ?
	fileData, err := GetFDFileById(fileId)
	if err != nil {
		return err
	}
	file, err := os.Open(fileData.FilePath)
	if err != nil {
		return err
	}
	r := csv.NewReader(file)

	if _, err := r.Read(); err != nil {
		return err
	}

	failedRows := []int64{}
	errorRows := []string{}
	records, err := r.ReadAll()
	if err != nil {
		return err
	}
	for i, record := range records {
		err := createFDHoldingFromRow(record, fileData.UserId)
		if err != nil {
			errorRows = append(errorRows, err.Error())
			failedRows = append(failedRows, int64(i))
		}
	}
	fileData.RowErrors = errorRows
	fileData.RowsFailed = failedRows
	fileData.Parsed = true
	db.Save(&fileData)
	return nil
}

func createFDHoldingFromRow(row []string, userId uint) error {
	accNo, ipFreq := row[0], row[2]
	amount, err := strconv.ParseFloat(row[1], 64)
	if err != nil {
		return err
	}
	ipRate, err := strconv.ParseFloat(row[3], 64)
	if err != nil {
		return err
	}
	mtAmount, err := strconv.ParseFloat(row[6], 64)
	if err != nil {
		return err
	}
	startDate, err := time.Parse(FileDtFormat, row[4])
	if err != nil {
		return err
	}
	endDate, err := time.Parse(FileDtFormat, row[5])
	if err != nil {
		return err
	}
	bankAccountUserMatch := account.CheckBankAccountAndUserId(userId, accNo)
	if !bankAccountUserMatch {
		return errors.New("bank account and users dont match")
	}
	fd := securities.FixedDeposit{
		Amount:    amount,
		IPFreq:    ipFreq,
		IPRate:    ipRate,
		MtAmount:  mtAmount,
		MtDate:    endDate,
		StartDate: startDate,
	}
	err = fd.Create()
	if err != nil {
		return err
	}
	ba, err := account.GetBankAccountByNumber(accNo)
	if err != nil {
		return err
	}
	fdh := FDHolding{
		FixedDeposit: fd,
		BankAccount:  ba,
	}
	fdh.create()
	return nil
}

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

func GetFDsByUser(userId uint) []FDHolding {
	var holdings []FDHolding
	dematIds, _ := account.GetBankAccountIdsByUser(userId)
	db.Joins("FixedDeposit").Find(&holdings, "bank_account_id IN ?", dematIds)
	return holdings
}
func (fdh FDHolding) getHoldings() HoldingSecurity {
	return HoldingSecurity{
		Name:         fmt.Sprintf("%f-%s-%s", fdh.FixedDeposit.Amount, fdh.FixedDeposit.IPFreq, fdh.FixedDeposit.MtDate),
		CurrentValue: fdh.FixedDeposit.Amount + securities.CalculateAccruedInterest(fdh.FixedDeposit),
		Invested:     fdh.FixedDeposit.Amount,
		Category:     "FD",
	}
}
