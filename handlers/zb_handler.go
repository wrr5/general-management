package handlers

import (
	"net/http"

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
