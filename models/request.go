package models

import (
	"gorm.io/gorm"
)

type Request struct {
	gorm.Model
	ResponseBody  string `gorm:"type:varchar(100);"`
	ResposeStatus int    `gorm:"type:integer;"`
	RequestBody   []byte `gorm:"type:varchar(100);"`
	Path          string `gorm:"type:varchar(100);"`
	Headers       []byte `gorm:"type:varchar(100);"`
	Method        string `gorm:"type:varchar(10);"`
	Host          string `gorm:"type:varchar(100);"`
	Ipadress      string `gorm:"type:varchar(100);"`
	TimeSpent	  string `gorm:"type:varchar(100);"` //or duration
}
