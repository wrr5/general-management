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
		// 设置商品路由
		setProductRoutes(api)
		// 设置工厂路由
		setFactoryRoutes(api)
		// 设置上传路由
		setUploadRoutes(api)
		// 设置售后路由
		setAfterSaleRoutes(api)
		// 设置物流轨迹路由
		setDeliveryRoutes(api)
		// 设置订单路由
		setOrderRoutes(api)
		// 设置快递路由
		setExpressRoutes(api)
		// 按照快递单号查询快递物流信息
		api.POST("delivery/query", handlers.GetLogisticsByNo)
	}

	// 404 处理
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "页面不存在"})
	})
	r.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", gin.H{}) })
	r.GET("/register", func(c *gin.Context) { c.HTML(http.StatusOK, "register.html", gin.H{}) })
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
		user.GET("/:id", middleware.RequireAuth(), handlers.GetUser)
		// 更新用户信息
		// user.PUT("/:id", handlers.UpdateUser)
		// 部分更新用户信息
		user.PATCH("/:id", handlers.PatchUser)
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
	zb.Use(middleware.RequireAuth())
	{
		zb.GET("", handlers.GetZb)
		zb.POST("", handlers.CreateZb)
		zb.DELETE("/:name", handlers.DeleteZb)
	}
}

func setUploadRoutes(r *gin.RouterGroup) {
	upload := r.Group("/upload")
	{
		upload.POST("shipments", middleware.RequireAuth(), handlers.UploadShipment)
	}
}

func setOrderRoutes(r *gin.RouterGroup) {
	order := r.Group("/orders")
	order.Use(middleware.RequireAuth())
	{
		order.POST("", handlers.CreateOrder)
		order.GET("", handlers.GetOrders)
		order.GET("/:id", handlers.GetOrder)
		order.PUT("/:id", handlers.UpdateOrder)
		order.PATCH("/:id", handlers.PatchOrder)
		order.DELETE("/:id", handlers.DeleteOrder)
	}
}

func setAfterSaleRoutes(r *gin.RouterGroup) {
	afterSale := r.Group("/aftersale")
	afterSale.Use(middleware.RequireAuth())
	{
		afterSale.POST("", handlers.CreateAfterSale)
		afterSale.GET("", handlers.GetAfterSale)
		afterSale.PATCH("", handlers.PatchAfterSale)
		afterSale.GET("/type", handlers.GetAfterSaleType)
		afterSale.GET("/request", handlers.GetAfterSaleRequest)
		afterSale.GET("/attribution", handlers.GetProblemAttribution)
	}
}

func setDeliveryRoutes(r *gin.RouterGroup) {
	delivery := r.Group("/delivery")
	delivery.Use(middleware.RequireAuth())
	{
		delivery.POST("", handlers.GetLogistics)
		delivery.POST("/accepted", handlers.GetAccepted)
		delivery.POST("/product-name", handlers.GetDeliveryByProductName)
		delivery.POST("/product", handlers.GetDeliveryByProductId)
	}
}

func setExpressRoutes(r *gin.RouterGroup) {
	express := r.Group("/express")
	express.Use(middleware.RequireAuth())
	{
		express.PATCH("", handlers.PatchExpress)
	}
}

func setProductRoutes(r *gin.RouterGroup) {
	product := r.Group("/products")
	product.Use(middleware.RequireAuth())
	{
		product.POST("", handlers.CreateProduct)       // 创建商品
		product.GET("", handlers.GetProducts)          // 获取商品列表
		product.GET("/:id", handlers.GetProduct)       // 获取单个商品
		product.PUT("/:id", handlers.UpdateProduct)    // 全量更新商品
		product.PATCH("/:id", handlers.PatchProduct)   // 部分更新商品
		product.DELETE("/:id", handlers.DeleteProduct) // 删除商品
	}
}

func setFactoryRoutes(r *gin.RouterGroup) {
	factory := r.Group("/factries")
	factory.Use(middleware.RequireAuth())
	{
		factory.POST("", handlers.CreateFactory)
		factory.GET("", handlers.GetFactories)
		factory.GET("/:id", handlers.GetFactory)
		factory.PUT("/:id", handlers.UpdateFactory)
		factory.PATCH("/:id", handlers.PatchFactory)
		factory.DELETE("/:id", handlers.DeleteFactory)
	}
}
