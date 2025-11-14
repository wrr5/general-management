package models

import (
	"time"
)

// Product 商品表模型
type Product struct {
	VzProductID       string    `gorm:"primaryKey;column:vz_product_id;type:varchar(50);not null;comment:商品ID"`
	ProductName       string    `gorm:"column:product_name;type:varchar(200);not null;index:idx_product_name;comment:商品名称"`
	Address           string    `gorm:"column:address;type:varchar(500);not null;comment:地址"`
	Specifications    *string   `gorm:"column:specifications;type:varchar(500);comment:可选规格"`
	Combo             *int      `gorm:"column:combo;type:int;index:idx_package;comment:套餐"`
	LivePrice         *float64  `gorm:"column:live_price;type:decimal(15,2);comment:直播价"`
	SupplyPrice       *float64  `gorm:"column:supply_price;type:decimal(15,2);comment:供货价"`
	FactoryEmployeeID *int      `gorm:"column:factory_employee_id;type:int;index:idx_factory_employee;comment:工厂员工ID"`
	CreatedTime       time.Time `gorm:"column:created_time;autoCreateTime;index:idx_created_time;comment:创建时间"`
	UpdatedTime       time.Time `gorm:"column:updated_time;autoUpdateTime;comment:更新时间"`

	// 关联关系
	FactoryEmployee *User `gorm:"foreignKey:FactoryEmployeeID;references:UserID"`
}

// TableName 设置表名
func (Product) TableName() string {
	return "products"
}
