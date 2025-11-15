package models

import (
	"time"
)

// Order 订单表模型
type Order struct {
	VzOrderID      int       `gorm:"primaryKey;column:vz_order_id;autoIncrement;comment:微赞订单ID"`
	VzProductID    string    `gorm:"column:vz_product_id;type:varchar(50);not null;index:idx_vz_product;comment:商品ID"`
	VzStoreID      string    `gorm:"column:vz_store_id;type:varchar(50);not null;index:idx_vz_store;comment:门店ID"`
	Quantity       int       `gorm:"column:quantity;type:int;not null;comment:数量"`
	Unit           int       `gorm:"column:unit;type:int;not null;comment:计件单位"`
	Specification  *string   `gorm:"column:specification;type:varchar(500);comment:规格"`
	CreatedTime    time.Time `gorm:"column:created_time;autoCreateTime;index:idx_created_time;comment:创建时间"`
	UpdatedTime    time.Time `gorm:"column:updated_time;autoUpdateTime;comment:更新时间"`
	ShipmentNumber string    `gorm:"column:shipment_number;not null;size:50" json:"shipment_number"`

	// 关联关系
	Product  *Product `gorm:"foreignKey:VzProductID;references:VzProductID"`
	Shipment Shipment `gorm:"foreignKey:ShipmentNumber;references:ShipmentNumber"`
	// 反向引用：一个订单对应多个快递记录
	Expresses []Express `gorm:"foreignKey:VzOrderID;references:VzOrderID"`
	// Store   *Store   `gorm:"foreignKey:VzStoreID;references:VzStoreID"`
}

// TableName 设置表名
func (Order) TableName() string {
	return "orders"
}
