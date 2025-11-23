package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wrr5/order-manage/global"
	"github.com/wrr5/order-manage/models"
)

func GetZb(c *gin.Context) {
	db := global.DB
	var zbs []models.Zb
	if result := db.Find(&zbs).Error; result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "查询直播间信息失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "查询直播间信息成功",
		"data":    zbs,
	})
}

func CreateZb(c *gin.Context) {
	type createRequest struct {
		ZbName string `form:"zb_name" json:"zb_name" binding:"required"`
		ZbID   string `form:"zb_id" json:"zb_id" binding:"required"`
	}
	db := global.DB
	var req createRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误: " + err.Error(),
		})
		return
	}
	zb := models.Zb{ZbName: req.ZbName, ZbID: req.ZbID}
	if result := db.Create(&zb); result.Error != nil {
		if strings.Contains(result.Error.Error(), "Duplicate entry") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "记录已存在",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "创建失败: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "创建直播间成功",
	})
}

func DeleteZb(c *gin.Context) {
	// 获取直播间name
	zbName := c.Param("name")
	db := global.DB
	result := db.Delete(&models.Zb{}, "zb_name = ?", zbName)

	if result.Error != nil {
		log.Printf("删除直播间失败: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "删除失败: " + result.Error.Error(),
		})
		return
	}

	// 检查是否实际删除了记录
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "直播间不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除直播间成功",
	})
}
