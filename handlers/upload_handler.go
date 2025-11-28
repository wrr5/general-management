package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wrr5/order-manage/global"
	"github.com/wrr5/order-manage/models"
	"gorm.io/gorm"
)

type ShipmentRequest struct {
	UserID         string         `json:"user_id" binding:"required"`
	ShipmentNumber string         `json:"shipmentnumber" binding:"required"`
	ProductID      string         `json:"productId" binding:"required"`
	ShipmentData   []ShipmentItem `json:"shipmentData" binding:"required,min=1"`
}

type ShipmentItem struct {
	Quantity       int     `json:"quantity" binding:"required,min=1"`
	Unit           int     `json:"unit" binding:"required"`
	StoreID        string  `json:"storeId" binding:"required"`
	Specification  *string `json:"specification"`
	TrackingNumber string  `json:"trackingNumber" binding:"required"`
}

func UploadShipment(c *gin.Context) {

	db := global.DB

	var req ShipmentRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误: " + err.Error(),
		})
		return
	}

	var user models.User
	if err := db.Where("user_id = ?", req.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的用户",
		})
		return
	}

	if user.UserType != 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "当前用户无权限",
		})
		return
	}

	var shipment models.Shipment
	if result := db.Where("shipment_number = ?", req.ShipmentNumber).First(&shipment); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "当前发货单已上传",
		})
		return
	}
	shipment.ShipmentNumber = req.ShipmentNumber
	shipment.ZbName = *user.ZbName

	// 使用事务处理所有数据库操作
	err := db.Transaction(func(tx *gorm.DB) error {
		shipment_num := models.Shipment{
			ShipmentNumber: req.ShipmentNumber,
			ZbName:         *user.ZbName,
		}
		if result := tx.Create(&shipment_num); result.Error != nil {
			return fmt.Errorf("创建发货单 %s 失败: %w", shipment.ShipmentNumber, result.Error)
		}
		for _, shipment := range req.ShipmentData {
			// 创建订单记录
			order := models.Order{
				VzProductID:    req.ProductID,
				VzStoreID:      shipment.StoreID,
				Quantity:       shipment.Quantity,
				Unit:           shipment.Unit,
				Specification:  shipment.Specification,
				ShipmentNumber: shipment_num.ShipmentNumber,
			}

			if result := tx.Create(&order); result.Error != nil {
				return fmt.Errorf("创建订单 %d 失败: %w", order.VzOrderID, result.Error)
			}

			// 处理快递单号
			trackingNumbers := strings.Split(shipment.TrackingNumber, "/")

			for _, trackingNum := range trackingNumbers {
				express := models.Express{
					VzOrderID:     order.VzOrderID,
					ExpressNumber: strings.TrimSpace(trackingNum),
					// 如果有其他字段，请继续添加
				}

				if result := tx.Create(&express); result.Error != nil {
					return fmt.Errorf("新增快递 %s 失败: %w", trackingNum, result.Error)
				}
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "上传发货单成功",
	})
}
