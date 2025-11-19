package models

import (
	"time"
)

// Express 快递表模型
type Express struct {
	ExpressNumber    string    `gorm:"primaryKey;column:express_number;type:varchar(100);not null;comment:快递单号"`
	VzOrderID        int       `gorm:"primaryKey;column:vz_order_id;not null;index:idx_vz_order_id;comment:订单微赞ID"`
	ActualQuantity   int       `gorm:"column:actual_quantity;not null;default:0;comment:实到数量"`
	CreatedTime      time.Time `gorm:"column:created_time;autoCreateTime;index:idx_created_time;comment:创建时间"`
	UpdatedTime      time.Time `gorm:"column:updated_time;autoUpdateTime;comment:更新时间"`
	IsAccepted       int8      `gorm:"column:is_accepted;type:tinyint;default:0;index:idx_is_accepted;comment:是否验收（0未验收，1已验收）"`
	UnboxingVideoURL *string   `gorm:"column:unboxing_video_url;type:varchar(500);comment:拆箱视频链接"`
	ActionUserID     *int      `gorm:"column:action_user_id;index:FK_action_user_id;comment:操作人ID"`

	// 关联关系
	Order      *Order `gorm:"foreignKey:VzOrderID;references:VzOrderID"`
	ActionUser *User  `gorm:"foreignKey:ActionUserID;references:UserID"`
}

// TableName 设置表名
func (Express) TableName() string {
	return "express"
}

// 状态常量
const (
	ExpressStatusNotAccepted = 0 // 未验收
	ExpressStatusAccepted    = 1 // 已验收
)
