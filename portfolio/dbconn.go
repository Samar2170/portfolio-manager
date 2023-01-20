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
	db.AutoMigrate(&MFTrade{})
	db.AutoMigrate(&MFHolding{})

	db.AutoMigrate(&FDHolding{})
	db.AutoMigrate(&FDFile{})
}

func handleError(err error) {
	if err != nil {
		log.Println(err)
	}
}
