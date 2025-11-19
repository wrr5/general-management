package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wrr5/order-manage/global"
	"github.com/wrr5/order-manage/models"
	"gorm.io/gorm"
)

func PatchExpress(c *gin.Context) {
	db := global.DB
	type ExpressUpdateReq struct {
		ExpressNumber    string   `json:"express_number" binding:"required"` // 快递单号（必填）
		VzOrderID        int      `json:"vz_order_id" binding:"required"`    // 微赞订单ID（必填）
		UserID           int      `json:"user_id" binding:"required"`        // 用户Id
		ActualQuantity   *float64 `json:"actual_quantity"`
		UnboxingVideoURL *string  `json:"unboxing_video_url"`
		IsAccepted       *int8    `json:"is_accepted"`
	}
	var req ExpressUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误: " + err.Error(),
		})
		return
	}
	// 查找现有的快递
	var existingExpress models.Express
	result := db.Where("express_number = ? AND vz_order_id = ?", req.ExpressNumber, req.VzOrderID).First(&existingExpress)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "快递不存在",
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

	if existingExpress.ActionUserID == nil {
		updates["action_user_id"] = &req.UserID
	} else {
		if *existingExpress.ActionUserID != req.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "当前用户无权限",
			})
			return
		}
	}

	// 处理可选字段更新
	if req.ActualQuantity != nil {
		updates["actual_quantity"] = *req.ActualQuantity
	}
	if req.IsAccepted != nil {
		updates["is_accepted"] = *req.IsAccepted
	}
	if req.UnboxingVideoURL != nil {
		updates["unboxing_video_url"] = *req.UnboxingVideoURL
	}

	// 如果没有要更新的字段，直接返回
	if len(updates) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "没有要更新的字段",
			"data":    existingExpress,
		})
		return
	}

	// 执行更新
	result = db.Model(&models.Express{}).
		Where("express_number = ? AND vz_order_id = ?", req.ExpressNumber, req.VzOrderID).
		Updates(updates)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新失败: " + result.Error.Error(),
		})
		return
	}

	// 重新查询更新后的数据
	var updatedExpress models.Express
	db.Where("express_number = ? AND vz_order_id = ?", req.ExpressNumber, req.VzOrderID).First(&updatedExpress)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新成功",
		"data":    updatedExpress,
	})
}
