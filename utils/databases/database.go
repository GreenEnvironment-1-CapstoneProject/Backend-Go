package databases

import (
	"fmt"
	"greenenvironment/configs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(c configs.GEConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DB_USER,
		c.DB_PASSWORD,
		c.DB_HOST,
		c.DB_PORT,
		c.DB_NAME,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		fmt.Println("terjadi kesalahan pada database, error:", err.Error())
		return nil, err
	}

	return db, err
}
