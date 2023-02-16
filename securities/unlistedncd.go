package securities

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type UnlistedNCD struct {
	*gorm.Model
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"uniqueIndex"`
	Symbol        string `gorm:"uniqueIndex"`
	IPFreq        string
	IPRate        float64
	IPDate        time.Time
	MtDate        time.Time
	FaceValue     float64
	MaturityValue float64
}

func (ulncd UnlistedNCD) GetOrCreate() UnlistedNCD {
	var nulncd UnlistedNCD
	db.FirstOrCreate(&nulncd, ulncd)
	return nulncd
}
func GetUnlistedNCDBySymbol(symbol string) (UnlistedNCD, error) {
	var ulncd UnlistedNCD
	symbol = strings.ToUpper(symbol)
	err := db.Where("symbol = ?", symbol).First(&ulncd).Error
	return ulncd, err
}
func CreateUnlistedNCD(name, symbol, ipRate, ipFreq, ipDate, mtDate, faceValue, mtValue string) (UnlistedNCD, error) {
	symbol = strings.ToUpper(symbol)
	if mtValue == "" {
		mtValue = faceValue
	}
	faceValueFloat, err := strconv.ParseFloat(faceValue, 64)
	if err != nil {
		return UnlistedNCD{}, fmt.Errorf("face_value should be float")
	}
	ipRateFloat, err := strconv.ParseFloat(ipRate, 64)
	if err != nil {
		return UnlistedNCD{}, fmt.Errorf("ip_rate should be float")
	}
	if ipRateFloat < 0 || ipRateFloat > 0.99 {
		return UnlistedNCD{}, fmt.Errorf("ip not in valid range")
	}
	if _, ok := ValidIPFreqs[ipFreq]; !ok {
		return UnlistedNCD{}, errors.New("ip freq not valid, should be A, M, MT, Q, SA")
	}
	mtValueFloat, err := strconv.ParseFloat(mtValue, 64)
	if err != nil {
		return UnlistedNCD{}, fmt.Errorf("maturity_value should be float")
	}
	ipDateParsed, err := time.Parse(DtFormat, ipDate)
	if err != nil {
		return UnlistedNCD{}, errors.New("ip_date not valid")
	}
	mtDateParsed, err := time.Parse(DtFormat, mtDate)
	if err != nil {
		return UnlistedNCD{}, errors.New("maturity_date not valid")
	}
	ncd := UnlistedNCD{
		Name:          name,
		Symbol:        symbol,
		IPRate:        ipRateFloat,
		IPFreq:        ipFreq,
		IPDate:        ipDateParsed,
		MtDate:        mtDateParsed,
		FaceValue:     faceValueFloat,
		MaturityValue: mtValueFloat,
	}
	ulncd := ncd.GetOrCreate()
	return ulncd, nil
}
