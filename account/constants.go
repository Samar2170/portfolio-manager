package account

import (
	"fmt"

	"github.com/spf13/viper"
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

const (
	DtFormat = "2006-01-02"
)

func init() {
	loadConfigFile()
	connect()
}
