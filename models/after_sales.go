package models

import (
	"time"
)

// AfterSale 售后表模型
type AfterSale struct {
	ExpressNumber          string    `gorm:"primaryKey;column:express_number;type:varchar(100);not null;comment:快递单号"`
	VzOrderID              int       `gorm:"primaryKey;column:vz_order_id;not null;index:FK_after_sales_vz_order_id_express_number;comment:订单微赞ID"`
	AfterSalesTypeID       int       `gorm:"column:after_sales_type_id;not null;index:idx_after_sales_type;comment:售后类型ID"`
	AfterSalesRequestID    int       `gorm:"column:after_sales_request_id;not null;index:idx_after_sales_request;comment:售后诉求ID"`
	ProblemAttributionID   int       `gorm:"column:problem_attribution_id;not null;index:FK_problem_attribution_id;comment:问题归属ID"`
	ProblemQuantity        int       `gorm:"column:problem_quantity;not null;comment:问题商品数量"`
	FreightVoucherURL      *string   `gorm:"column:freight_voucher_url;type:varchar(500);comment:运费凭证图片链接"`
	FreightAmount          *float64  `gorm:"column:freight_amount;type:decimal(10,2);comment:运费金额"`
	StoreExpressNumber     *string   `gorm:"column:store_express_number;type:varchar(100);index:idx_store_express;comment:门店快递单号"`
	FactoryExpressNumber   *string   `gorm:"column:factory_express_number;type:varchar(100);index:idx_factory_express;comment:厂家快递单号"`
	ProcessStatus          int8      `gorm:"column:process_status;type:tinyint;not null;default:1;index:idx_process_status;comment:处理进度（1待处理，2处理中，3门店待发，4厂家待收，5厂家待发，6门店待收，7运费待结算，8已完结，9已驳回）"`
	InitiatorID            int       `gorm:"column:initiator_id;not null;index:idx_initiator;comment:发起人ID"`
	CreatedTime            time.Time `gorm:"column:created_time;autoCreateTime;index:idx_created_time;comment:创建时间"`
	ProgressUpdatedTime    time.Time `gorm:"column:progress_updated_time;autoCreateTime;index:idx_progress_time;comment:进度更新时间"`
	StoreRemarks           *string   `gorm:"column:store_remarks;type:varchar(255);comment:门店备注"`
	FactoryRemarks         *string   `gorm:"column:factory_remarks;type:varchar(255);comment:厂家备注"`
	CustomerServiceRemarks *string   `gorm:"column:customer_service_remarks;type:varchar(255);comment:客服备注"`
	FactoryRequestID       *int      `gorm:"column:factory_request_id;index:FK_factory_request_id;comment:经工厂确认的诉求ID"`

	// 关联关系
	Express            *Express            `gorm:"foreignKey:ExpressNumber,VzOrderID;references:ExpressNumber,VzOrderID"`
	AfterSalesType     *AfterSalesType     `gorm:"foreignKey:AfterSalesTypeID;references:ID"`
	AfterSalesRequest  *AfterSalesRequest  `gorm:"foreignKey:AfterSalesRequestID;references:ID"`
	ProblemAttribution *ProblemAttribution `gorm:"foreignKey:ProblemAttributionID;references:ID"`
	Initiator          *User               `gorm:"foreignKey:InitiatorID;references:UserID"`
	FactoryRequest     *AfterSalesRequest  `gorm:"foreignKey:FactoryRequestID;references:ID"`
}

// TableName 设置表名
func (AfterSale) TableName() string {
	return "after_sales"
}

// 处理进度状态常量
const (
	ProcessStatusPending          int8 = 1 // 待处理
	ProcessStatusProcessing       int8 = 2 // 处理中
	ProcessStatusStoreToSend      int8 = 3 // 门店待发
	ProcessStatusFactoryToReceive int8 = 4 // 厂家待收
	ProcessStatusFactoryToSend    int8 = 5 // 厂家待发
	ProcessStatusStoreToReceive   int8 = 6 // 门店待收
	ProcessStatusFreightPending   int8 = 7 // 运费待结算
	ProcessStatusCompleted        int8 = 8 // 已完结
	ProcessStatusRejected         int8 = 9 // 已驳回
)

// AfterSalesRequest 售后诉求表模型
type AfterSalesRequest struct {
	ID                 int       `gorm:"primaryKey;column:id;autoIncrement;comment:主键ID"`
	RequestDescription string    `gorm:"column:request_description;type:varchar(500);not null;uniqueIndex:idx_request_description;comment:售后诉求描述"`
	CreatedTime        time.Time `gorm:"column:created_time;autoCreateTime;comment:创建时间"`
	UpdatedTime        time.Time `gorm:"column:updated_time;autoUpdateTime;comment:更新时间"`
	Status             int8      `gorm:"column:status;type:tinyint;default:1;index:idx_status;comment:状态（1启用，0停用）"`
}

// TableName 设置表名
func (AfterSalesRequest) TableName() string {
	return "after_sales_request"
}

// AfterSalesType 售后类型表模型
type AfterSalesType struct {
	ID              int       `gorm:"primaryKey;column:id;autoIncrement;comment:主键ID"`
	TypeDescription string    `gorm:"column:type_description;type:varchar(200);not null;uniqueIndex:idx_type_description;comment:售后类型描述"`
	CreatedTime     time.Time `gorm:"column:created_time;autoCreateTime;comment:创建时间"`
	UpdatedTime     time.Time `gorm:"column:updated_time;autoUpdateTime;comment:更新时间"`
	Status          int8      `gorm:"column:status;type:tinyint;default:1;index:idx_status;comment:状态（1启用，0停用）"`
}

// TableName 设置表名
func (AfterSalesType) TableName() string {
	return "after_sales_type"
}

// 状态常量
const (
	AfterSalesTypeStatusActive   = 1 // 启用
	AfterSalesTypeStatusInactive = 0 // 停用
)

// ProblemAttribution 问题归属表模型
type ProblemAttribution struct {
	ID                            int    `gorm:"primaryKey;column:id;autoIncrement;comment:主键ID"`
	ProblemAttributionDescription string `gorm:"column:problem_attribution_description;type:varchar(255);comment:问题归属描述"`
}

// TableName 设置表名
func (ProblemAttribution) TableName() string {
	return "problem_attribution"
}
