package models

import (
	"gorm.io/gorm"
)

// DeliveryOrder 发货单模型
type DeliveryOrder struct {
	gorm.Model
	DeliveryOrderNo   string `gorm:"uniqueIndex;type:varchar(50)"` // 发货单号，设为唯一索引
	StoreID           string `gorm:"type:varchar(20);index"`       // 门店ID
	StoreName         string `gorm:"type:varchar(100)"`            // 门店名称
	Period            string `gorm:"type:varchar(20);index"`       // 期数
	ProductID         string `gorm:"type:varchar(20);index"`       // 商品ID
	ProductName       string `gorm:"type:varchar(200)"`            // 商品名称
	Specification     string `gorm:"type:varchar(100)"`            // 规格
	Supplier          string `gorm:"type:varchar(100)"`            // 供应商
	DeliveryCount     int    `gorm:"column:delivery_count"`        // 实发单数
	DeliveryQuantity  int    `gorm:"column:delivery_quantity"`     // 实发商品数
	DeliveryInteger   int    `gorm:"column:delivery_integer"`      // 实发整数
	DeliveryScattered int    `gorm:"column:delivery_scattered"`    // 实发零散
	Status            string `gorm:"type:varchar(20);index"`       // 状态
	Receiver          string `gorm:"type:varchar(50)"`             // 收货人
	Phone             string `gorm:"type:varchar(20)"`             // 电话
	Address           string `gorm:"type:varchar(500)"`            // 地址
	LogisticsCompany  string `gorm:"type:varchar(100)"`            // 物流公司
	LogisticsPhone    string `gorm:"type:varchar(20)"`             // 联系电话
	LogisticsAddress  string `gorm:"type:varchar(500)"`            // 联系地址
	TrackingNumber    string `gorm:"type:varchar(500)"`            // 物流单号
}

// TableName 指定表名
func (DeliveryOrder) TableName() string {
	return "delivery_orders"
}
