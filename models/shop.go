package models

import (
	"gorm.io/gorm"
)

type Shop struct {
	gorm.Model
	Name    string `gorm:"type:varchar(50);not null" json:"name"`
	Address string `gorm:"type:varchar(255)" json:"address"`
}
