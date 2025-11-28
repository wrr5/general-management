package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wrr5/order-manage/global"
	"github.com/wrr5/order-manage/models"
	"gorm.io/gorm"
)

func CreateFactory(c *gin.Context) {
	type createRequest struct {
		VzFactoryID string `form:"vz_factory_id" json:"vz_factory_id" binding:"required"`
		FactoryName string `form:"factory_name" json:"factory_name" binding:"required"`
	}
	db := global.DB
	var req createRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误: " + err.Error(),
		})
		return
	}
	factory := models.Factory{VzFactoryID: req.VzFactoryID, FactoryName: req.FactoryName}
	if result := db.Create(&factory); result.Error != nil {
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
		"message": "创建工厂成功",
	})
}

func GetFactories(c *gin.Context) {
	db := global.DB
	var factories []models.Factory
	if result := db.Find(&factories).Error; result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "查询工厂列表失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "查询工厂列表成功",
		"data":    factories,
	})
}

func GetFactory(c *gin.Context) {
	id := c.Param("id")
	db := global.DB
	var factory models.Factory
	if result := db.Where("vz_factory_id = ?", id).First(&factory); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "查询工厂信息失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "查询工厂信息成功",
		"data":    factory,
	})
}

func UpdateFactory(c *gin.Context) {

}

func PatchFactory(c *gin.Context) {

}

func DeleteFactory(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "工厂ID不能为空",
		})
		return
	}

	// 检查工厂是否存在
	var existingFactory models.Factory
	if err := global.DB.First(&existingFactory, "vz_factory_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "工厂不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	// 可选：检查是否有关联的产品（避免外键约束错误）
	var productCount int64
	global.DB.Model(&models.Product{}).Where("factory_employee_id IN (SELECT user_id FROM users WHERE vz_factory_id = ?)", id).Count(&productCount)
	if productCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "该工厂下有关联的产品，无法删除",
			"data": gin.H{
				"related_products_count": productCount,
			},
		})
		return
	}

	// 执行删除
	result := global.DB.Where("vz_factory_id = ?", id).Delete(&models.Factory{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除失败: " + result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "工厂不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "工厂删除成功",
	})
}
