package models

import (
	"gorm.io/gorm"
)

type Request struct {
	gorm.Model
	ResponseBody  string  `gorm:"type:text;" json:"responsebody"`
	ResposeStatus int     `gorm:"type:integer;" json:"code"`
	RequestBody   string  `gorm:"type:text;" json:"requestbody"`
	Path          string  `gorm:"type:text;" json:"path"`
	Headers       string  `gorm:"type:text;" json:"headers"`
	Method        string  `gorm:"type:varchar(10);" json:"method"`
	Host          string  `gorm:"type:varchar(100);" json:"host"`
	Ipadress      string  `gorm:"type:varchar(100);" json:"ipaddress"`
	TimeSpent     float64 `gorm:"type:float;" json:"timespent"`
}
