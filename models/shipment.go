package models

type Shipment struct {
	ShipmentNumber string `gorm:"primaryKey;column:shipment_number;size:50" json:"shipment_number"`
	ZbName         string `gorm:"primaryKey;column:zb_name;size:50" json:"zb_name"`

	// 关联关系
	Zb Zb `gorm:"foreignKey:ZbName;references:ZbName"`

	// 反向关联
	// Orders []Order `gorm:"foreignKey:ShipmentNumber;references:ShipmentNumber"`
}

// 表名设置
func (Shipment) TableName() string {
	return "shipment"
}
