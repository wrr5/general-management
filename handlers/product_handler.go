package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wrr5/order-manage/global"
	"github.com/wrr5/order-manage/models"
	"gorm.io/gorm"
)

func CreateProduct(c *gin.Context) {
	type createProductRequest struct {
		VzProductID       string   `json:"vz_product_id" binding:"required"`
		ProductName       string   `json:"product_name" binding:"required"`
		Address           string   `json:"address" binding:"required"`
		Specifications    *string  `json:"specifications,omitempty"`
		Combo             *int     `json:"combo,omitempty"`
		LivePrice         *float64 `json:"live_price,omitempty"`
		SupplyPrice       *float64 `json:"supply_price,omitempty"`
		FactoryEmployeeID *int     `json:"factory_employee_id,omitempty"`
	}
	var req createProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}
	// 转换为模型
	product := models.Product{
		VzProductID:       req.VzProductID,
		ProductName:       req.ProductName,
		Address:           req.Address,
		Specifications:    req.Specifications,
		Combo:             req.Combo,
		LivePrice:         req.LivePrice,
		SupplyPrice:       req.SupplyPrice,
		FactoryEmployeeID: req.FactoryEmployeeID,
	}

	// 保存到数据库
	if err := global.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建商品失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "商品创建成功",
	})
}

func GetProducts(c *gin.Context) {
	db := global.DB
	var products []models.Product
	if result := db.Find(&products).Error; result != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "查询商品列表失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "查询商品列表成功",
		"data":    products,
	})
}

func GetProduct(c *gin.Context) {
	id := c.Param("id")
	db := global.DB
	var product models.Product
	if result := db.Where("vz_product_id = ?", id).First(&product); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "查询商品信息失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "查询商品信息成功",
		"data":    product,
	})
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	// 参数验证
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "商品ID不能为空",
		})
		return
	}

	var req struct {
		ProductName       string   `json:"product_name" binding:"required"`
		Address           string   `json:"address" binding:"required"`
		Specifications    *string  `json:"specifications,omitempty"`
		Combo             *int     `json:"combo,omitempty"`
		LivePrice         *float64 `json:"live_price,omitempty"`
		SupplyPrice       *float64 `json:"supply_price,omitempty"`
		FactoryEmployeeID *int     `json:"factory_employee_id,omitempty"`
	}

	// 绑定请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 检查商品是否存在
	var existingProduct models.Product
	if err := global.DB.First(&existingProduct, "vz_product_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "商品不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	// 执行更新
	updateData := models.Product{
		ProductName:       req.ProductName,
		Address:           req.Address,
		Specifications:    req.Specifications,
		Combo:             req.Combo,
		LivePrice:         req.LivePrice,
		SupplyPrice:       req.SupplyPrice,
		FactoryEmployeeID: req.FactoryEmployeeID,
	}

	result := global.DB.Model(&models.Product{}).
		Where("vz_product_id = ?", id).
		Updates(updateData)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新失败: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "商品更新成功",
	})
}

func PatchProduct(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "商品ID不能为空",
		})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 检查商品是否存在
	var existingProduct models.Product
	if err := global.DB.First(&existingProduct, "vz_product_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "商品不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	// 移除不能更新的字段
	delete(req, "vz_product_id")
	delete(req, "created_time")
	delete(req, "updated_time")

	// 如果没有任何可更新字段
	if len(req) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "没有提供可更新的字段",
		})
		return
	}

	// 执行部分更新
	result := global.DB.Model(&models.Product{}).
		Where("vz_product_id = ?", id).
		Updates(req)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新失败: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "商品部分更新成功",
	})
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "商品ID不能为空",
		})
		return
	}

	// 检查商品是否存在
	var existingProduct models.Product
	if err := global.DB.First(&existingProduct, "vz_product_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "商品不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	// 执行删除
	result := global.DB.Where("vz_product_id = ?", id).Delete(&models.Product{})

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
			"message": "商品不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "商品删除成功",
	})
}
