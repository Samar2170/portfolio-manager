package securities

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Samar2170/portfolio-manager/utils"
	"gorm.io/gorm"
)

type ListedNCD struct {
	*gorm.Model
	ID              uint   `gorm:"primaryKey"`
	Name            string `gorm:"uniqueIndex"`
	Symbol          string `gorm:"uniqueIndex"`
	SecurityCode    string `gorm:"unique"`
	Exchange        string
	IPFreq          string
	IPRate          float64
	IPDate          time.Time
	MtDate          time.Time
	FaceValue       float64
	MaturityValue   float64
	NextIPDate      time.Time
	AccruedInterest float64
	IssueDate       time.Time
}

func (lncd ListedNCD) GetOrCreate() ListedNCD {
	var nlncd ListedNCD
	db.FirstOrCreate(&nlncd, lncd)
	return nlncd
}

func GetListedNCDBySymbol(symbol string) (ListedNCD, error) {
	var lncd ListedNCD
	symbol = strings.ToUpper(symbol)
	err := db.Where("symbol = ?", symbol).First(&lncd).Error
	return lncd, err
}

func getListedNCDsByIds(ids []uint) ([]ListedNCD, error) {
	var lncds []ListedNCD
	err := db.Where("id IN ?", ids).Find(&lncds).Error
	return lncds, err
}

func GetNCDList() ([]string, error) {
	var ncds []ListedNCD
	var symbols []string

	err := db.Find(&ncds).Error
	if err != nil {
		return symbols, err
	}
	for _, ncd := range ncds {
		symbols = append(symbols, ncd.Symbol)
	}
	return symbols, nil
}
func CreateListedNCD(name, symbol, securityCode, exchange, ipRate, ipFreq, ipDate, mtDate, faceValue, mtValue string) (ListedNCD, error) {
	symbol = strings.ToUpper(symbol)
	if mtValue == "" {
		mtValue = faceValue
	}
	faceValueFloat, err := strconv.ParseFloat(faceValue, 64)
	if err != nil {
		return ListedNCD{}, fmt.Errorf("face_value should be float")
	}
	ipRateFloat, err := strconv.ParseFloat(ipRate, 64)
	if err != nil {
		return ListedNCD{}, fmt.Errorf("ip_rate should be float")
	}
	if ipRateFloat < 0 || ipRateFloat > 0.99 {
		return ListedNCD{}, fmt.Errorf("ip not in valid range")
	}
	if _, ok := ValidIPFreqs[ipFreq]; !ok {
		return ListedNCD{}, errors.New("ip freq not valid, should be A, M, MT, Q, SA")
	}
	mtValueFloat, err := strconv.ParseFloat(mtValue, 64)
	if err != nil {
		return ListedNCD{}, fmt.Errorf("maturity_value should be float")
	}
	ipDateParsed, err := time.Parse(DtFormat, ipDate)
	if err != nil {
		return ListedNCD{}, errors.New("ip_date not valid")
	}
	mtDateParsed, err := time.Parse(DtFormat, mtDate)
	if err != nil {
		return ListedNCD{}, errors.New("maturity_date not valid")
	}
	ncd := ListedNCD{
		Name:          name,
		Symbol:        symbol,
		Exchange:      exchange,
		SecurityCode:  securityCode,
		IPRate:        ipRateFloat,
		IPDate:        ipDateParsed,
		IPFreq:        ipFreq,
		MtDate:        mtDateParsed,
		FaceValue:     faceValueFloat,
		MaturityValue: mtValueFloat,
	}
	lncd := ncd.GetOrCreate()
	return lncd, nil
}

func CalculateNextIPDateListedNCD(lncd ListedNCD) time.Time {
	var nextIPDate time.Time
	today := time.Now()
	ipdate := lncd.IPDate
	var currIpDate time.Time
	if ipdate.Before(today) {
		switch lncd.IPFreq {
		case "A":
			currIpDate = time.Date(today.Year(), ipdate.Month(), ipdate.Day(), 0, 0, 0, 0, time.Local)
			if currIpDate.Before(today) {
				nextIPDate = currIpDate.AddDate(1, 0, 0)
			} else {
				nextIPDate = currIpDate
			}

		case "M":
			currIpDate = time.Date(today.Year(), today.Month(), ipdate.Day(), 0, 0, 0, 0, time.Local)
			nextIPDate = currIpDate.AddDate(0, 1, 0)
		case "Q":
			nextIPDate = utils.GetNextQuarter(today)
		case "MT":
			nextIPDate = lncd.MtDate
		case "SA":
			nextIPDate = utils.GetNextHY(today)

		}
	} else {
		nextIPDate = currIpDate
	}
	return nextIPDate
}

func CalculateAccruedInterestListedNCD(lncd ListedNCD) float64 {
	t := time.Now()
	if lncd.IPFreq == "MT" {
		totalInterest := lncd.MaturityValue - lncd.FaceValue
		totalTime := lncd.MtDate.Sub(lncd.IssueDate) / 24
		dailyInterest := totalInterest / float64(totalTime)
		timeDone := t.Sub(lncd.IssueDate) / 24
		return float64(timeDone) * dailyInterest
	} else {
		var daysDiff int
		var intMultiple float64
		annualInterest := lncd.IPRate * lncd.FaceValue
		switch lncd.IPFreq {
		case "M":
			daysDiff = int(t.Day())
			intMultiple = float64(daysDiff) / 30 * 0.0833

		case "Q":
			cqDate := utils.GetCurrentQuarterFirstDate(t)
			daysDiff = int(t.Sub(cqDate).Hours() / 24)
			intMultiple = float64(daysDiff) / 90 * 0.25
		case "SA":
			cqDate := utils.GetCurrentHYFirstDate(t)
			daysDiff = int(t.Sub(cqDate).Hours() / 24)
			intMultiple = float64(daysDiff) / 180 * 0.5
		case "A":
			ipDate := time.Date(t.Year(), lncd.IssueDate.Month(), lncd.IssueDate.Day(), 0, 0, 0, 0, time.UTC)
			if t.After(ipDate) {
				delta := ipDate.Sub(t)
				daysDiff = int(delta.Hours() / 24)
				intMultiple = float64(daysDiff) / 365
			} else {
				ipDate = time.Date(t.Year()-1, lncd.IssueDate.Month(), lncd.IssueDate.Day(), 0, 0, 0, 0, time.UTC)
				delta := t.Sub(ipDate)
				daysDiff := int(delta.Hours() / 24)
				intMultiple = float64(daysDiff) / 365
			}
		}
		accruedInt := annualInterest * intMultiple
		return accruedInt
	}
}

func CalculateAccruedInterestAllListedNCDs() error {
	var ids []uint
	db.Raw("SELECT id FROM listed_ncds").Scan(&ids)
	batches := len(ids) / BATCHSIZE

	for i := 0; i <= batches; i++ {
		newIds := ids[i*BATCHSIZE : i*BATCHSIZE+BATCHSIZE]
		lncds, err := getListedNCDsByIds(newIds)
		if err != nil {
			return err
		}
		for _, lncd := range lncds {
			accruedInterest := CalculateAccruedInterestListedNCD(lncd)
			db.Model(&lncd).Where("id = ?", lncd.ID).Update("accrued_interest", accruedInterest)

		}
	}
	return nil
}
