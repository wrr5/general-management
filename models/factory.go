package models

import (
	"time"
)

// Factory 工厂表模型
type Factory struct {
	VzFactoryID string    `gorm:"primaryKey;column:vz_factory_id;type:varchar(50);not null;comment:工厂微赞ID"`
	FactoryName string    `gorm:"column:factory_name;type:varchar(200);not null;index:idx_factory_name;comment:工厂名称"`
	CreatedTime time.Time `gorm:"column:created_time;autoCreateTime;index:idx_created_time;comment:创建时间"`
	UpdatedTime time.Time `gorm:"column:updated_time;autoUpdateTime;comment:更新时间"`
}

// TableName 设置表名
func (Factory) TableName() string {
	return "factories"
}
