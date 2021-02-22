package models

import "gorm.io/gorm"

type Request struct {
	gorm.Model
	ID            uint
	ResponseBody  string `gorm:"type:varchar(100);"`
	ResposeStatus int
	RequestBody   string `gorm:"type:varchar(100);"`
	Path          string `gorm:"type:varchar(100);"`
	Headers       string `gorm:"type:varchar(100);"`
	Method        string `gorm:"type:varchar(10);"`
	Host          string `gorm:"type:varchar(100);"`
}
