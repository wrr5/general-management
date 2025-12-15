package main

import (
	"html/template"
	"strings"

	"github.com/gin-gonic/gin"
	// "github.com/thinkerou/favicon"
	"github.com/wrr5/order-manage/config"
	"github.com/wrr5/order-manage/global"
	"github.com/wrr5/order-manage/router"
	"github.com/wrr5/order-manage/tools"
)

func main() {
	config.Init()
	if config.AppConfig.Server.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	tools.InitDB()
	// 初始化token管理器
	global.TM.StartAutoRefresh()

	r := router.SetupRouter()

	// 注册模板函数
	r.SetFuncMap(template.FuncMap{
		"splitTrackingNumber": func(trackingNumber string) []string {
			if trackingNumber == "" {
				return []string{}
			}
			// 分割单号，支持 / 分隔符
			numbers := strings.Split(trackingNumber, "/")
			// 清理空白字符
			result := make([]string, 0, len(numbers))
			for _, num := range numbers {
				trimmed := strings.TrimSpace(num)
				if trimmed != "" {
					result = append(result, trimmed)
				}
			}
			return result
		},
		"sub": func(a, b int) int { return a - b },
		"len": func(arr []string) int { return len(arr) },
		"lt":  func(a, b int) bool { return a < b },
	})

	r.LoadHTMLGlob("templates/*.html")
	r.Static("/static", "./static")
	r.Static("/uploads", "./uploads")
	// r.Use(favicon.New("./static/images/favicon.ico"))

	r.Run("0.0.0.0:" + config.AppConfig.Server.Port)
}
