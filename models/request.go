package models

import (
	"gorm.io/gorm"
)

type Request struct {
	gorm.Model
	ResponseBody  string  `gorm:"type:varchar(1000);"`
	ResposeStatus int     `gorm:"type:integer;"`
	RequestBody   string  `gorm:"type:varchar(1000);"`
	Path          string  `gorm:"type:varchar(100);"`
	Headers       string  `gorm:"type:varchar(1000);"`
	Method        string  `gorm:"type:varchar(10);"`
	Host          string  `gorm:"type:varchar(100);"`
	Ipadress      string  `gorm:"type:varchar(100);"`
	TimeSpent     float64 `gorm:"type:float;"`
}
