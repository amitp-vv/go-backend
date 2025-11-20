package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("realestate.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	// Auto-migrate all models
	err = DB.AutoMigrate(
		&User{},
		&Wallet{},
		&Property{},
		&Claim{},
		&Distribution{},
	)
	return err
}
