package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wrr5/order-manage/global"
	"github.com/wrr5/order-manage/models"
	"github.com/wrr5/order-manage/services"
)

func GetLogistics(c *gin.Context) {
	type queryExpressRequest struct {
		VzStoreID  string `json:"VzStoreID" binding:"required"`
		IsAccepted int    `json:"is_accepted"`
		StateText  string `json:"stateText"`
	}
	var req queryExpressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误: " + err.Error(),
		})
		return
	}

	type ExpressResponse struct {
		ExpressNumber string         `json:"ExpressNumber"`
		DeliveryTime  time.Time      `json:"deliveryTime"` // 发货时间
		LatestTrace   services.Trace `json:"latestTrace"`
	}
	type OrderResponse struct {
		VzOrderID       int               `json:"VzOrderID"`
		ExpressResponse []ExpressResponse `json:"deliveries"`
		ProductID       string            `json:"productID"`
		ProductName     string            `json:"productName"`
		IsAcceptedCount int               `json:"IsAcceptedCount"`
		ProductSum      int               `json:"ProductSum"`
	}

	var Orders []OrderResponse

	db := global.DB
	var orders []models.Order
	var express []models.Express
	db.Where("vz_store_id = ?", req.VzStoreID).Preload("Product").Find(&orders)
	for _, order := range orders {
		count := 0
		newOrder := OrderResponse{
			VzOrderID:   order.VzOrderID,
			ProductID:   order.VzProductID,
			ProductName: order.Product.ProductName,
			ProductSum:  order.Quantity * order.Unit,
		}
		db.Where("vz_order_id = ?", order.VzOrderID).Find(&express)
		for _, exp := range express {
			if exp.IsAccepted == int8(req.IsAccepted) {
				count += 1
				logisticsResponse, err := services.QueryDelivery(exp.ExpressNumber)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err.Error(),
					})
					return
				}
				if req.StateText == "" {
					newExp := ExpressResponse{
						ExpressNumber: exp.ExpressNumber,
						DeliveryTime:  exp.CreatedTime,
						LatestTrace:   logisticsResponse.DataObj.LogisticsInfo.Traces[0],
					}
					newOrder.ExpressResponse = append(newOrder.ExpressResponse, newExp)
				} else {
					if logisticsResponse.DataObj.LogisticsInfo.StateText == req.StateText {
						newExp := ExpressResponse{
							ExpressNumber: exp.ExpressNumber,
							DeliveryTime:  exp.CreatedTime,
							LatestTrace:   logisticsResponse.DataObj.LogisticsInfo.Traces[0],
						}
						newOrder.ExpressResponse = append(newOrder.ExpressResponse, newExp)
					}
				}
			}
		}
		newOrder.IsAcceptedCount = count
		Orders = append(Orders, newOrder)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "查询物流成功",
		"orders":  Orders,
	})
}
