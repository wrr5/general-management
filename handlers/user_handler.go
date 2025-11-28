package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wrr5/order-manage/global"
	"github.com/wrr5/order-manage/models"
	"github.com/wrr5/order-manage/services"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(c *gin.Context) {
	db := global.DB

	type createRequest struct {
		PhoneNumber string `form:"phone_mumber" json:"phone_mumber" binding:"required,len=11"`
		RealName    string `form:"real_name" json:"real_name" binding:"required"`
		UserType    string `form:"user_type" json:"user_type" binding:"required"`
		ZbName      string `form:"zb_name" json:"zb_name" binding:"required"`
		Password    string `form:"password" json:"password" binding:"required,min=6,max=20"`
		VzStoreID   string `form:"vz_store_id" json:"vz_store_id"`
		VzFactoryID string `form:"vz_factory_id" json:"vz_factory_id"`
	}

	var req createRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误: " + err.Error(),
		})
		return
	}

	var existingUser models.User
	if err := db.Where("phone_number = ? AND zb_name = ?", req.PhoneNumber, req.ZbName).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "用户已存在",
		})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 其他数据库错误
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "查询用户失败",
		})
		return
	}

	if req.UserType == "2" {
		_, err := services.ValidatePhoneName(req.PhoneNumber, req.RealName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "门店用户信息验证失败",
			})
			return
		}
	}
	value, err := strconv.ParseInt(req.UserType, 10, 8)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用户类型不合法",
		})
		return
	}
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "密码加密失败",
		})
		return
	}
	user := models.User{
		PhoneNumber: req.PhoneNumber,
		RealName:    req.RealName,
		UserType:    int8(value),
		ZbName:      &req.ZbName,
		Password:    string(hashedPassword),
	}

	if req.VzStoreID != "" {
		user.VzStoreID = &req.VzStoreID
	}

	if req.VzFactoryID != "" {
		user.VzFactoryID = &req.VzFactoryID
	}

	if result := db.Create(&user).Error; result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用户创建失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "创建用户成功",
		"user":    user,
	})
}

func GetUsers(c *gin.Context) {
	type getRequest struct {
		UserType string `form:"user_type" json:"user_type"`
	}
	var req getRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误: " + err.Error(),
		})
		return
	}
	db := global.DB
	var users []models.User
	if req.UserType != "" {
		if err := db.Where("user_type = ?", req.UserType).Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "查询用户失败",
			})
			return
		}
	} else {
		if err := db.Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "查询用户失败",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "查询用户列表成功",
		"data":    users,
		"total":   len(users),
	})
}

func GetUser(c *gin.Context) {
	// 获取字符串类型的 ID
	idStr := c.Param("id")

	// 转换为整数
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ID 必须是数字",
		})
		return
	}
	var user models.User
	db := global.DB
	if result := db.First(&user, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "未查询到用户",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "查询用户成功",
		"data":    user,
	})
}

func PatchUser(c *gin.Context) {
	db := global.DB
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID不能为空"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式错误"})
		return
	}

	// 检查用户是否存在
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 过滤允许更新的字段
	allowedFields := []string{"user_type", "real_name", "phone_number", "vz_store_id", "vz_factory_id", "zb_name"}
	filteredData := make(map[string]interface{})

	for _, field := range allowedFields {
		if value, exists := updateData[field]; exists {
			filteredData[field] = value
		}
	}

	// 特殊处理密码字段
	if password, exists := updateData["password"]; exists && password != "" {
		pwdStr, ok := password.(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "密码格式错误",
			})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwdStr), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "密码加密失败",
			})
			return
		}
		// 转换为字符串
		filteredData["password"] = string(hashedPassword)
	}

	if len(filteredData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有有效的更新字段"})
		return
	}

	// 执行更新
	if err := db.Model(&user).Updates(filteredData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新用户失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "用户更新成功",
	})
}
