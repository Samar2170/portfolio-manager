package securities

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
	db.AutoMigrate(&Stock{})
	// db.AutoMigrate(&ListedNCD{})
	// db.AutoMigrate(&Crypto{})
	db.AutoMigrate(&MutualFund{})
	// db.AutoMigrate(&UnlistedNCD{})
	db.AutoMigrate(&FixedDeposit{})
}

func handleError(err error) {
	if err != nil {
		log.Println(err)
	}
}
