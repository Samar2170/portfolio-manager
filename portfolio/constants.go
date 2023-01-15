package portfolio

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	DtFormat         = "2006-01-02"
	SmtpAddress      = "smtp.gmail.com"
	SmtpAddressWPort = "smtp.gmail.com:587"
)

var DBURI string
var EMAILID string
var EMAILPASSWORD string

func loadConfigFile() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	DBURI = fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable",
		viper.Get("DBHOST"),
		viper.Get("DBUSER"),
		viper.Get("DBNAME"),
		viper.Get("DBPASSWORD"),
	)
	EMAILID = viper.GetString("EMAIL_ID")
	EMAILPASSWORD = viper.GetString("EMAIL_PASSWORD")
}

func init() {
	loadConfigFile()
	connect()
}

var BankNames = map[string]struct{}{
	"KOTAK":    {},
	"HDFC":     {},
	"INDUSIND": {},
	"RBL":      {},
	"YES":      {},
	"SBI":      {},
	"AXIS":     {},
}

var ValidIPFreqs = map[string]struct{}{
	"A":  {},
	"M":  {},
	"MT": {},
	"Q":  {},
	"SA": {},
	// "QAD": {}, // quarterly adjusting, Bank format.
	// "SAD": {}, // semi annual adjusting, Govt format.
}
