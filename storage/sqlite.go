package storage

import (
	"github.com/ichtrojan/horus/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func Connect() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "horus.db")

	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.Request{}).Error; err != nil {
		return nil, err
	}

	return db, nil
}
