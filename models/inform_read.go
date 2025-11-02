package models

import (
	"time"

	"gorm.io/gorm"
)

// InformRead 通知读取记录表（多对多关系）
type InformRead struct {
	gorm.Model

	UserID   uint `gorm:"primaryKey;autoIncrement:false" json:"user_id"`
	InformID uint `gorm:"primaryKey;autoIncrement:false" json:"inform_id"`

	ReadAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"read_at"`

	// 外键关系
	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	Inform Inform `gorm:"foreignKey:InformID;constraint:OnDelete:CASCADE" json:"inform"`
}

// TableName 设置表名
func (InformRead) TableName() string {
	return "inform_reads"
}
