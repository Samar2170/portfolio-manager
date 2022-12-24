package account

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	*gorm.Model
	Username string `gorm:"type:varchar(100);uniqueIndex"`
	Password string `gorm:"type:varchar(100)"`
	Email    string `gorm:"type:varchar(100);uniqueIndex"`
}

type GeneralAccount struct {
	*gorm.Model
	Code   string `gorm:"unique_index"`
	User   User   `gorm:"foreignKey:UserId"`
	UserId uint
}

type DematAccount struct {
	*gorm.Model
	Code   string `gorm:"uniqueIndex"`
	User   User   `gorm:"foreignKey:UserId"`
	UserId uint
	Broker string
}

type BankAccount struct {
	*gorm.Model
	AccountNo string `gorm:"uniqueIndex"`
	User      User   `gorm:"foreignKey:UserId"`
	UserId    uint
	Bank      string
}

func (u *User) GetOrCreate() (User, error) {
	var existingUser User
	_ = db.Where("username = ?", u.Username).First(&existingUser).Error
	if existingUser.Username == u.Username {
		return User{}, errors.New("username already exists")
	}
	err := db.Create(&u).Error
	return *u, err
}

func GetUserById(userId uint) (User, error) {
	var user User
	err := db.First("id = ?", userId).Error
	return user, err
}

func GetUserByUsername(username string) (User, error) {
	var user User
	err := db.First("username = ", username).Error
	return user, err
}

func (ga *GeneralAccount) GetOrCreate() (GeneralAccount, error) {
	err := db.Create(&ga).Error
	return *ga, err
}

func (da *DematAccount) GetOrCreate() (DematAccount, error) {
	err := db.Create(&da).Error
	return *da, err
}
func (ba *BankAccount) GetOrCreate() (BankAccount, error) {
	err := db.Create(&ba).Error
	return *ba, err
}
