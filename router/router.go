package router

import (
	"html/template"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wrr5/order-manage/handlers"
	"github.com/wrr5/order-manage/middleware"
)

// SetupRouter 配置所有路由
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有来源
		AllowMethods:     []string{"*"}, // 允许所有方法
		AllowHeaders:     []string{"*"}, // 允许所有头
		ExposeHeaders:    []string{"*"}, // 暴露所有头
		AllowCredentials: true,          // 允许凭证
		MaxAge:           12 * 60 * 60,  // 预检请求缓存时间
	}))
	// 添加模板函数
	r.SetFuncMap(template.FuncMap{})

	// 创建 API 路由组
	api := r.Group("/api")
	{
		// 设置认证路由
		setAuthRoutes(api)
		// 设置用户路由
		setUserRoutes(api)
		// 设置直播路由
		setZbRoutes(api)
		// 设置上传路由
		setUploadRoutes(api)
		// 设置售后路由
		setAfterSaleRoutes(api)
		// 设置物流轨迹路由
		setDeliveryRoutes(api)
	}

	// 404 处理
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "页面不存在"})
	})
	r.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", gin.H{}) })

	return r
}

func setAuthRoutes(r *gin.RouterGroup) {

	r.POST("/login", handlers.Login)

}

func setUserRoutes(r *gin.RouterGroup) {
	user := r.Group("/users")
	{
		// 创建用户
		user.POST("", handlers.CreateUser)
		// 获取用户列表（带分页和筛选）
		user.GET("", middleware.RequireAuth(), handlers.GetUsers)
		// 获取单个用户详情
		user.GET("/:id", handlers.GetUser)
		// 更新用户信息
		// user.PUT("/:id", handlers.UpdateUser)
		// 部分更新用户信息
		// user.PATCH("/:id", handlers.PartialUpdateUser)
		// 删除用户
		// user.DELETE("/:id", handlers.DeleteUser)
		// 获取当前登录用户信息
		// user.GET("/me", handlers.GetCurrentUser)
		// 更新当前用户密码
		// user.PUT("/me/password", handlers.UpdatePassword)
	}
}

func setZbRoutes(r *gin.RouterGroup) {
	zb := r.Group("/zb")
	{
		zb.GET("", handlers.GetZb)
	}
}

func setUploadRoutes(r *gin.RouterGroup) {
	upload := r.Group("/upload")
	{
		upload.POST("shipments", middleware.RequireAuth(), handlers.UploadShipment)
	}
}

func setAfterSaleRoutes(r *gin.RouterGroup) {
	afterSale := r.Group("/aftersale")
	afterSale.Use(middleware.RequireAuth())
	{
		afterSale.POST("", handlers.CreateAfterSale)
		afterSale.GET("", handlers.GetAfterSale)
		afterSale.PATCH("", handlers.PartialUpdateAfterSale)
	}
}

func setDeliveryRoutes(r *gin.RouterGroup) {
	delivery := r.Group("/delivery")
	delivery.Use(middleware.RequireAuth())
	{
		delivery.POST("", handlers.GetLogistics)
	}
}
