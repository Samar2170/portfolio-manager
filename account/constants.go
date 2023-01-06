package account

import (
	"fmt"

	"github.com/spf13/viper"
)

var DBURI string

func loadConfigFile() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	DBURI = fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable",
		viper.Get("DBHOST"),
		viper.Get("DBUSER"),
		viper.Get("DBNAME"),
		viper.Get("DBPASSWORD"),
	)
}

const (
	DtFormat = "2006-01-02"
)

func init() {
	loadConfigFile()
	connect()
}
