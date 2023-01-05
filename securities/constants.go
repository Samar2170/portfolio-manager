package securities

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	NiftyStocksFile  = "./Assets/MCAP31032022.xlsx"
	BSEStocksFile    = "./Assets/BSEEquities.xlsx"
	MFSchemesFile    = "./Assets/SchemeData2312221939SS.xlsx"
	NSE500StocksFile = "./Assets/ind_nifty500list.xlsx"
)

var ValidIPFreqs = map[string]struct{}{
	"A":  {},
	"M":  {},
	"MT": {},
	"Q":  {},
	"SA": {},
}
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
