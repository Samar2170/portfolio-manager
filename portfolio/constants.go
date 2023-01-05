package portfolio

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	DtFormat = "2006-01-02"
)

var DBURI string

func loadConfigFile() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	fmt.Println(viper.Get("DBPORT"))
	fmt.Println(viper.Get("DBHOST"))
	fmt.Println(viper.Get("DBUSER"))
	fmt.Println(viper.Get("DBPASSWORD"))
	fmt.Println(viper.Get("DBNAME"))

	DBURI = fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable",
		viper.Get("DBHOST"),
		viper.Get("DBUSER"),
		viper.Get("DBNAME"),
		viper.Get("DBPASSWORD"),
	)
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
}
