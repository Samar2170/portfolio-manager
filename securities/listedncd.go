package securities

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

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
