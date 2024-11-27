package databases

import (
	DataAdmin "greenenvironment/features/admin/repository"
	DataUser "greenenvironment/features/users/repository"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	db.AutoMigrate(&DataUser.User{})
	db.AutoMigrate(&DataAdmin.Admin{})
	return nil
}
