package models

// Zb 直播间表模型
type Zb struct {
	ZbName string `gorm:"primaryKey;column:zb_name;type:varchar(50);not null;comment:直播间名称"`
	ZbID   string `gorm:"column:zbId;type:varchar(100);not null;comment:直播ID"`

	// 反向关联
	// Shipments []Shipment `gorm:"foreignKey:ZbName;references:ZbName"`
}

// TableName 设置表名
func (Zb) TableName() string {
	return "zb"
}
