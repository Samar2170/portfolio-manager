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
	BATCHSIZE        = 10
	DtFormat         = "2006-01-02"
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
