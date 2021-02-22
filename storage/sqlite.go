package storage

import (
	"github.com/ichtrojan/horus/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("horus.db"), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(models.Request{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
