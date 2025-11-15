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
		VzStoreID string `json:"VzStoreID" binding:"required"`
		StateText string `json:"stateText"`
	}
	var req queryExpressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误: " + err.Error(),
		})
		return
	}

	type ExpressResponse struct {
		ExpressNumber  string         `json:"ExpressNumber"`
		ActualQuantity float64        `json:"ActualQuantity"`
		DeliveryTime   time.Time      `json:"deliveryTime"` // 发货时间
		LatestTrace    services.Trace `json:"latestTrace"`
	}
	type OrderResponse struct {
		VzOrderID         int               `json:"VzOrderID"`
		ExpressResponse   []ExpressResponse `json:"deliveries"`
		ProductID         string            `json:"productID"`
		ProductName       string            `json:"productName"`
		Accepted          int               `json:"Accepted"`
		UnAccepted        int               `json:"UnAccepted"`
		ActualQuantitySum float64           `json:"ActualQuantitySum"`
		ProductSum        int               `json:"ProductSum"`
	}

	var Orders []OrderResponse

	db := global.DB
	var orders []models.Order
	var express []models.Express
	unAcceptedtotle := 0
	totle := 0
	db.Where("vz_store_id = ?", req.VzStoreID).Preload("Product").Find(&orders)
	for _, order := range orders {
		newOrder := OrderResponse{
			VzOrderID:   order.VzOrderID,
			ProductID:   order.VzProductID,
			ProductName: order.Product.ProductName,
			ProductSum:  order.Quantity * order.Unit,
		}
		db.Where("vz_order_id = ?", order.VzOrderID).Find(&express)
		actualQuantitySum := 0.0
		accepted := 0
		unAccepted := 0
		for _, exp := range express {
			actualQuantitySum += exp.ActualQuantity
			switch exp.IsAccepted {
			case 0:
				unAccepted += 1
			case 1:
				accepted += 1
			}
			// time.Sleep(800 * time.Millisecond)
			logisticsResponse, err := services.QueryDelivery(exp.ExpressNumber)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			if req.StateText == "" {
				newExp := ExpressResponse{
					ExpressNumber:  exp.ExpressNumber,
					ActualQuantity: exp.ActualQuantity,
					DeliveryTime:   exp.CreatedTime,
					LatestTrace:    logisticsResponse.DataObj.LogisticsInfo.Traces[0],
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
		newOrder.ActualQuantitySum = actualQuantitySum
		newOrder.Accepted = accepted
		newOrder.UnAccepted = unAccepted
		unAcceptedtotle += unAccepted
		totle += unAccepted + accepted
		Orders = append(Orders, newOrder)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"message":         "查询物流成功",
		"totle":           totle,
		"unAcceptedtotle": unAcceptedtotle,
		"data":            Orders,
	})
}

func GetProductExpress(c *gin.Context) {
	type queryExpressRequest struct {
		VzStoreID   string `json:"VzStoreID" binding:"required"`
		ProductName string `json:"ProductName"`
	}
	var req queryExpressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误: " + err.Error(),
		})
		return
	}

	type ExpressResponse struct {
		ExpressNumber string           `json:"ExpressNumber"`
		DeliveryTime  time.Time        `json:"deliveryTime"` // 发货时间
		Trace         []services.Trace `json:"latestTrace"`
	}
	type OrderResponse struct {
		VzOrderID       int               `json:"VzOrderID"`
		ExpressResponse []ExpressResponse `json:"deliveries"`
		ProductID       string            `json:"productID"`
		ProductName     string            `json:"productName"`
		ProductSum      int               `json:"ProductSum"`
	}

	var respOrders []OrderResponse
	db := global.DB
	var orders []models.Order
	if req.ProductName != "" {
		var productIDs []string
		db.Model(&models.Product{}).Where("product_name LIKE ?", "%"+req.ProductName+"%").Pluck("vz_product_id", &productIDs)
		db.Where("vz_store_id = ? AND vz_product_id IN (?)", req.VzStoreID, productIDs).
			Preload("Product").Preload("Expresses").Find(&orders)

		for _, order := range orders {
			newOrder := OrderResponse{
				VzOrderID:   order.VzOrderID,
				ProductID:   order.VzProductID,
				ProductName: order.Product.ProductName,
				ProductSum:  order.Quantity * order.Unit,
			}

			for _, exp := range order.Expresses {
				logisticsResponse, err := services.QueryDelivery(exp.ExpressNumber)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err.Error(),
					})
					return
				}

				newExp := ExpressResponse{
					ExpressNumber: exp.ExpressNumber,
					DeliveryTime:  exp.CreatedTime,
					Trace:         logisticsResponse.DataObj.LogisticsInfo.Traces,
				}
				newOrder.ExpressResponse = append(newOrder.ExpressResponse, newExp)
			}
			respOrders = append(respOrders, newOrder)
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "按商品名筛选成功",
			"data":    respOrders,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "未传递商品名",
		"data":    respOrders,
	})
}
