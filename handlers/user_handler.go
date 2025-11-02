package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wrr5/general-management/global"
	"github.com/wrr5/general-management/models"
)

func ShowUserPage(c *gin.Context) {
	userInterface, _ := c.Get("user")
	user := userInterface.(models.User)
	db := global.DB

	var users []models.User
	err := db.Order("created_at DESC").Find(&users).Error
	if err != nil {
		// 处理错误
		c.HTML(http.StatusOK, "user.html", gin.H{
			"CurrentPath": c.Request.URL.Path,
			"err":         err.Error(),
			"users":       []models.User{},
		})
		return
	}
	c.HTML(http.StatusOK, "user.html", gin.H{
		"CurrentPath": c.Request.URL.Path,
		"user":        user,
		"users":       users,
	})
}
