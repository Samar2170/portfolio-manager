package account

import "gorm.io/gorm"

type JobRecord struct {
	*gorm.Model
	Name    string
	Fields  string
	Args    string
	Success bool
	Error   string
}

func (jr JobRecord) Create() error {
	err := db.Create(&jr).Error
	return err
}
