package portfolio

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func connect() {
	var err error
	db, err = gorm.Open(postgres.Open(DBURI), &gorm.Config{})
	handleError(err)
	db.AutoMigrate(&StockTrade{})
	db.AutoMigrate(&StockHolding{})
}

func handleError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func init() {
	connect()
}
