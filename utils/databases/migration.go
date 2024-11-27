package databases

import (
	DataUser "greenenvironment/features/users/repository"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	db.AutoMigrate(&DataUser.User{})
	return nil
}
