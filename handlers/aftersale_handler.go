package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wrr5/order-manage/global"
	"github.com/wrr5/order-manage/models"
	"gorm.io/gorm"
)

func GetAfterSaleType(c *gin.Context) {
	db := global.DB
	var afterSaleTypes []models.AfterSalesType
	if result := db.Find(&afterSaleTypes); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "查询售后类型失败：" + result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "查询售后类型成功",
		"data":    afterSaleTypes,
	})
}

func GetAfterSaleRequest(c *gin.Context) {
	db := global.DB
	var afterSaleRequests []models.AfterSalesRequest
	if result := db.Find(&afterSaleRequests); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "查询售后诉求失败：" + result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "查询售后诉求成功",
		"data":    afterSaleRequests,
	})
}

func CreateAfterSale(c *gin.Context) {
	type AfterSaleCreateReq struct {
		ExpressNumber        string  `json:"express_number" binding:"required"`         // 快递单号（必填）
		VzOrderID            int     `json:"vz_order_id" binding:"required"`            // 微赞订单ID（必填）
		AfterSalesTypeID     int     `json:"after_sales_type_id" binding:"required"`    // 售后类型ID（必填）
		AfterSalesRequestID  int     `json:"after_sales_request_id" binding:"required"` // 售后诉求ID（必填）
		ProblemQuantity      int     `json:"problem_quantity" binding:"required,min=1"` // 问题商品数量（必填）
		InitiatorID          int     `json:"initiator_id" binding:"required"`           // 发起人ID（必填）
		ProblemAttribution   string  `json:"problem_attribution"`                       // 问题归属（可选）
		UnboxingVideoURL     string  `json:"unboxing_video_url"`                        // 拆箱视频链接（可选）
		FreightVoucherURL    string  `json:"freight_voucher_url"`                       // 运费凭证图片链接（可选）
		FreightAmount        float64 `json:"freight_amount"`                            // 运费金额（可选）
		StoreExpressNumber   string  `json:"store_express_number"`                      // 门店快递单号（可选）
		FactoryExpressNumber string  `json:"factory_express_number"`                    // 厂家快递单号（可选）
		ProcessStatus        int8    `json:"process_status"`                            // 处理进度（可选，默认1）
	}

	var req AfterSaleCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误: " + err.Error(),
		})
		return
	}

	// 验证快递单号和订单号有效性
	var express models.Express
	if err := global.DB.Where("express_number = ? and vz_order_id = ?", req.ExpressNumber, req.VzOrderID).First(&express).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "无效的快递单号和订单号",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "服务器错误",
		})
		return
	}

	// 验证售后类型是否存在
	var afterSalesType models.AfterSalesType
	if err := global.DB.Where("id = ?", req.AfterSalesTypeID).First(&afterSalesType).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "无效的售后类型",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "服务器错误",
		})
		return
	}

	// 验证售后诉求是否存在
	var afterSalesRequest models.AfterSalesRequest
	if err := global.DB.Where("id = ?", req.AfterSalesRequestID).First(&afterSalesRequest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "无效的售后诉求",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "服务器错误",
		})
		return
	}

	// 验证发起人是否存在
	var user models.User
	if err := global.DB.Where("user_id = ?", req.InitiatorID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "发起人不存在",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "服务器错误",
		})
		return
	}

	// 创建售后单记录
	afterSale := models.AfterSale{
		ExpressNumber:       req.ExpressNumber,
		VzOrderID:           req.VzOrderID,
		AfterSalesTypeID:    req.AfterSalesTypeID,
		AfterSalesRequestID: req.AfterSalesRequestID,
		ProblemQuantity:     req.ProblemQuantity,
		InitiatorID:         req.InitiatorID,
	}
	// 处理字符串可选字段 - 只有非空时才赋值
	if req.ProblemAttribution != "" {
		afterSale.ProblemAttribution = &req.ProblemAttribution
	}
	if req.UnboxingVideoURL != "" {
		afterSale.UnboxingVideoURL = &req.UnboxingVideoURL
	}
	if req.FreightVoucherURL != "" {
		afterSale.FreightVoucherURL = &req.FreightVoucherURL
	}
	if req.StoreExpressNumber != "" {
		afterSale.StoreExpressNumber = &req.StoreExpressNumber
	}
	if req.FactoryExpressNumber != "" {
		afterSale.FactoryExpressNumber = &req.FactoryExpressNumber
	}

	// 处理数值可选字段 - 只有非零时才赋值
	if req.FreightAmount != 0 {
		afterSale.FreightAmount = &req.FreightAmount
	}
	if req.ProcessStatus != 0 {
		afterSale.ProcessStatus = req.ProcessStatus
	}

	// 保存到数据库
	if err := global.DB.Create(&afterSale).Error; err != nil {
		// 处理唯一约束冲突
		if strings.Contains(err.Error(), "Duplicate entry") {
			c.JSON(http.StatusConflict, gin.H{
				"error": "售后单已存在",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建售后单失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "创建售后单成功",
		"data":    afterSale,
	})
}

func GetAfterSale(c *gin.Context) {
	expressNumber := c.Query("express_number")
	vzOrderID := c.Query("vz_order_id")

	// 如果提供了联合主键，查询特定记录
	db := global.DB
	var existingAfterSale models.AfterSale
	if expressNumber != "" && vzOrderID != "" {
		// 按联合主键查询单一记录
		result := db.Where("express_number = ? AND vz_order_id = ?", expressNumber, vzOrderID).First(&existingAfterSale)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "未找到售后单",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "查询售后单成功",
			"data":    existingAfterSale,
		})
		return
	}

	// 如果没有提供主键，返回列表或分页结果
	// page := c.DefaultQuery("page", "1")
	// limit := c.DefaultQuery("limit", "20")
	var afterSales []models.AfterSale
	if result := db.Find(&afterSales); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "售后单查询错误：" + result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "查询售后单列表成功",
		"data":    afterSales,
	})
}

// AfterSaleUpdateReq 更新请求结构体
type AfterSaleUpdateReq struct {
	ExpressNumber        string   `json:"express_number" binding:"required"` // 快递单号（必填）
	VzOrderID            int      `json:"vz_order_id" binding:"required"`    // 微赞订单ID（必填）
	ProblemAttribution   *string  `json:"problem_attribution"`               // 问题归属（可选）
	UnboxingVideoURL     *string  `json:"unboxing_video_url"`                // 拆箱视频链接（可选）
	ProblemQuantity      *int     `json:"problem_quantity"`                  // 问题商品数量（可选）
	FreightVoucherURL    *string  `json:"freight_voucher_url"`               // 运费凭证图片链接（可选）
	FreightAmount        *float64 `json:"freight_amount"`                    // 运费金额（可选）
	StoreExpressNumber   *string  `json:"store_express_number"`              // 门店快递单号（可选）
	FactoryExpressNumber *string  `json:"factory_express_number"`            // 厂家快递单号（可选）
	ProcessStatus        *int8    `json:"process_status"`                    // 处理进度（可选）
	InitiatorID          *int     `json:"initiator_id"`                      // 发起人ID（可选）
}

func PatchAfterSale(c *gin.Context) {
	db := global.DB
	var req AfterSaleUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误: " + err.Error(),
		})
		return
	}

	// 查找现有的售后单
	var existingAfterSale models.AfterSale
	result := db.Where("express_number = ? AND vz_order_id = ?", req.ExpressNumber, req.VzOrderID).First(&existingAfterSale)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "售后单不存在",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "查询失败: " + result.Error.Error(),
			})
		}
		return
	}

	// 构建更新字段映射
	updates := make(map[string]interface{})

	// 处理可选字段更新
	if req.ProblemAttribution != nil {
		updates["problem_attribution"] = *req.ProblemAttribution
	}
	if req.UnboxingVideoURL != nil {
		updates["unboxing_video_url"] = *req.UnboxingVideoURL
	}
	if req.ProblemQuantity != nil {
		updates["problem_quantity"] = *req.ProblemQuantity
	}
	if req.FreightVoucherURL != nil {
		updates["freight_voucher_url"] = *req.FreightVoucherURL
	}
	if req.FreightAmount != nil {
		updates["freight_amount"] = *req.FreightAmount
	}
	if req.StoreExpressNumber != nil {
		updates["store_express_number"] = *req.StoreExpressNumber
	}
	if req.FactoryExpressNumber != nil {
		updates["factory_express_number"] = *req.FactoryExpressNumber
	}
	if req.ProcessStatus != nil {
		updates["process_status"] = *req.ProcessStatus
		// 如果更新了处理进度，同时更新进度更新时间
		updates["progress_updated_time"] = time.Now()
	}
	if req.InitiatorID != nil {
		updates["initiator_id"] = *req.InitiatorID
	}

	// 如果没有要更新的字段，直接返回
	if len(updates) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "没有要更新的字段",
			"data":    existingAfterSale,
		})
		return
	}

	// 执行更新
	result = db.Model(&models.AfterSale{}).
		Where("express_number = ? AND vz_order_id = ?", req.ExpressNumber, req.VzOrderID).
		Updates(updates)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新失败: " + result.Error.Error(),
		})
		return
	}

	// 重新查询更新后的数据
	var updatedAfterSale models.AfterSale
	db.Where("express_number = ? AND vz_order_id = ?", req.ExpressNumber, req.VzOrderID).First(&updatedAfterSale)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新成功",
		"data":    updatedAfterSale,
	})
}
