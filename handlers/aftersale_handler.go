package handlers

import "github.com/gin-gonic/gin"

func CreateAfterSale(c *gin.Context) {
	type AfterSaleCreateReq struct {
		ExpressNumber        string   `json:"express_number" binding:"required"`         // 快递单号（必填）
		VzOrderID            string   `json:"vz_order_id" binding:"required"`            // 订单微赞ID（必填）
		AfterSalesTypeID     int      `json:"after_sales_type_id" binding:"required"`    // 售后类型ID（必填）
		AfterSalesRequestID  int      `json:"after_sales_request_id" binding:"required"` // 售后诉求ID（必填）
		ProblemAttribution   *string  `json:"problem_attribution"`                       // 问题归属（可选）
		UnboxingVideoURL     *string  `json:"unboxing_video_url"`                        // 拆箱视频链接（可选）
		ProblemQuantity      int      `json:"problem_quantity" binding:"required,min=1"` // 问题商品数量（必填）
		FreightVoucherURL    *string  `json:"freight_voucher_url"`                       // 运费凭证图片链接（可选）
		FreightAmount        *float64 `json:"freight_amount"`                            // 运费金额（可选）
		StoreExpressNumber   *string  `json:"store_express_number"`                      // 门店快递单号（可选）
		FactoryExpressNumber *string  `json:"factory_express_number"`                    // 厂家快递单号（可选）
		ProcessStatus        int8     `json:"process_status"`                            // 处理进度（可选，默认1）
		InitiatorID          int      `json:"initiator_id" binding:"required"`           // 发起人ID（必填）
	}

}

func GetAfterSale(c *gin.Context) {

}

func PartialUpdateAfterSale(c *gin.Context) {

}
