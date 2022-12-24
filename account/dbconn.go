package account

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
	db.AutoMigrate(&User{})
	db.AutoMigrate(&GeneralAccount{})
	db.AutoMigrate(&DematAccount{})
	db.AutoMigrate(&BankAccount{})
}

func handleError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func init() {
	connect()
}
