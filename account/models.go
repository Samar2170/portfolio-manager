package account

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	*gorm.Model
	ID       uint
	Username string `gorm:"type:varchar(100);uniqueIndex"`
	Password string `gorm:"type:varchar(100)"`
	Email    string `gorm:"type:varchar(100);uniqueIndex"`
}

type UserAccountStatus struct {
	*gorm.Model
	ID               uint
	UserId           uint
	User             User `gorm:"foreignKey:UserId"`
	GaAccountCreated bool `gorm:"default:false"`
}

type GeneralAccount struct {
	*gorm.Model
	ID     uint
	Code   string `gorm:"unique_index"`
	User   User   `gorm:"foreignKey:UserId"`
	UserId uint
}

type DematAccount struct {
	*gorm.Model
	ID     uint
	Code   string `gorm:"uniqueIndex"`
	User   User   `gorm:"foreignKey:UserId"`
	UserId uint
	Broker string
}

type BankAccount struct {
	*gorm.Model
	ID        uint
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
	err := db.Where("id = ?", userId).First(&user).Error
	return user, err
}

func GetUserByUsername(username string) (User, error) {
	var user User
	err := db.Where("username = ?", username).First(&user).Error
	return user, err
}

func (ga *GeneralAccount) Create() (GeneralAccount, error) {
	err := db.Create(&ga).Error
	return *ga, err
}

func (da *DematAccount) Create() (DematAccount, error) {
	err := db.Create(&da).Error
	return *da, err
}
func (ba *BankAccount) Create() (BankAccount, error) {
	err := db.Create(&ba).Error
	return *ba, err
}

func GetDematAccountByCode(code string) (DematAccount, error) {
	var da DematAccount
	err := db.Where("code = ?", code).First(&da).Error
	return da, err
}

func (uas UserAccountStatus) GetOrCreate() UserAccountStatus {
	var nuas UserAccountStatus
	_ = db.FirstOrCreate(&nuas, uas)
	return nuas
}

func GetBankAccountByNumber(accNumber string) (BankAccount, error) {
	var ba BankAccount
	err := db.First(&ba, "account_no = ?", accNumber).Error
	return ba, err
}
func GetBankAccountByUserNBank(userId uint, bankName string) (BankAccount, error) {
	var ba BankAccount
	err := db.Where("user_id = ? AND bank = ?", userId, bankName).Error
	return ba, err
}

func CheckBankAccountAndUserId(userId uint, accountNumber string) bool {
	var record BankAccount
	err := db.Where("account_no = ? AND user_id = ?", accountNumber, userId).First(&record).Error
	return err != gorm.ErrRecordNotFound
}
func CheckDematAccountAndUserId(userId uint, dematAccCode string) bool {
	var record DematAccount
	err := db.Where("account_no = ? AND user_id = ?", dematAccCode, userId).First(&record).Error
	return err != gorm.ErrRecordNotFound
}

func GetBankAccountsByUser(userId uint) ([]BankAccount, error) {
	var accounts []BankAccount
	err := db.Where("user_id = ?", userId).Find(&accounts).Error
	return accounts, err
}
func GetDematAccountsByUser(userId uint) ([]DematAccount, error) {
	var accounts []DematAccount
	err := db.Where("user_id = ?", userId).Find(&accounts).Error
	return accounts, err
}
