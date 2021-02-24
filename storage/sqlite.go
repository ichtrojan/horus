package storage

import (
	"github.com/ichtrojan/horus/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error

func Connect() (*gorm.DB, error) {

	db, err = gorm.Open("sqlite3", "horus.db")
	//db, err := gorm.Open(sqlite.Open("horus.db"), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Request{})

	if err != nil {
		return nil, err
	}

	return db, nil
}

func DB() *gorm.DB {
	return db
}
