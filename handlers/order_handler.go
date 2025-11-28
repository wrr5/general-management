package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wrr5/order-manage/global"
	"github.com/wrr5/order-manage/models"
)

func CreateOrder(c *gin.Context) {
	type OrderRequest struct {
		Quantity       int     `json:"quantity" binding:"required,min=1"`
		Unit           int     `json:"unit" binding:"required"`
		ProductID      string  `json:"productId" binding:"required"`
		StoreID        string  `json:"storeId" binding:"required"`
		Specification  *string `json:"specification"`
		TrackingNumber string  `json:"trackingNumber" binding:"required"`
		ShipmentNumber string  `json:"shipmentnumber" binding:"required"`
	}
	var req OrderRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误: " + err.Error(),
		})
		return
	}
	db := global.DB
	// 创建订单记录
	order := models.Order{
		VzProductID:    req.ProductID,
		VzStoreID:      req.StoreID,
		Quantity:       req.Quantity,
		Unit:           req.Unit,
		Specification:  req.Specification,
		ShipmentNumber: req.ShipmentNumber,
	}

	if result := db.Create(&order); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建订单失败: " + result.Error.Error(),
		})
		return
	}

	// 处理快递单号
	trackingNumbers := strings.Split(req.TrackingNumber, "/")

	for _, trackingNum := range trackingNumbers {
		express := models.Express{
			VzOrderID:     order.VzOrderID,
			ExpressNumber: strings.TrimSpace(trackingNum),
			// 如果有其他字段，请继续添加
		}

		if result := db.Create(&express); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "创建快递" + trackingNum + "失败: " + result.Error.Error(),
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "创建订单成功",
	})
}

func GetOrders(c *gin.Context) {

}

func GetOrder(c *gin.Context) {

}

func UpdateOrder(c *gin.Context) {

}

func PatchOrder(c *gin.Context) {

}

func DeleteOrder(c *gin.Context) {

}
