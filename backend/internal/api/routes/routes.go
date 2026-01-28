package routes

import (
	"github.com/gin-gonic/gin"
	"online-mall/internal/api/controller"
	"online-mall/internal/api/middleware"
)

// SetupRoutes 设置路由
func SetupRoutes() *gin.Engine {
	// 设置运行模式
	if !gin.IsDebugging() {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建引擎
	r := gin.Default()

	// 使用中间件
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.Recovery())
	r.Use(middleware.RequestID())

	// 健康检查
	r.GET("/health", healthCheck)

	// API路由组
	api := r.Group("/api")
	{
		// 认证相关路由（不需要JWT）
		auth := api.Group("/auth")
		{
			auth.POST("/login", controller.Login)
			auth.POST("/register", controller.Register)
			auth.POST("/refresh-token", controller.RefreshToken)
			auth.POST("/logout", middleware.JWTAuth(), controller.Logout)
		}

		// 用户相关路由（需要JWT）
		user := api.Group("/users")
		user.Use(middleware.JWTAuth())
		{
			// 用户信息
			user.GET("/profile", controller.GetUserInfo)
			user.PUT("/profile", controller.UpdateUserInfo)
			user.PUT("/password", controller.UpdatePassword)

			// 管理员路由 - 待实现
			/*
				admin := user.Group("")
				admin.Use(middleware.RequireAdmin())
				{
					admin.GET("", controller.GetUserList)
					admin.GET("/:id", controller.GetUserDetail)
					admin.PUT("/:id", controller.UpdateUser)
					admin.PUT("/:id/status", controller.UpdateUserStatus)
					admin.DELETE("/:id", controller.DeleteUser)
				}
			*/
		}

		// 地址管理路由 - 待实现
		/*
			addresses := api.Group("/addresses")
			addresses.Use(middleware.JWTAuth())
			{
				addresses.GET("", controller.GetAddressList)
				addresses.POST("", controller.CreateAddress)
				addresses.GET("/:id", controller.GetAddressDetail)
				addresses.PUT("/:id", controller.UpdateAddress)
				addresses.DELETE("/:id", controller.DeleteAddress)
				addresses.PUT("/:id/default", controller.SetDefaultAddress)
			}
		*/

		// 商品相关路由
		products := api.Group("/products")
		{
			products.GET("", controller.GetProducts)
			products.GET("/:id", controller.GetProduct)
			products.GET("/:id/skus", controller.GetProductSkus)
			products.GET("/hot", controller.GetHotProducts)
			products.GET("/new", controller.GetNewProducts)

			// 管理员路由
			adminProducts := products.Group("")
			adminProducts.Use(middleware.RequireAdmin())
			{
				adminProducts.POST("", controller.CreateProduct)
				adminProducts.PUT("/:id", controller.UpdateProduct)
				adminProducts.DELETE("/:id", controller.DeleteProduct)
				adminProducts.PUT("/:id/status", controller.UpdateProductStatus)
			}
		}

		// 商品分类路由
		categories := api.Group("/categories")
		{
			categories.GET("", controller.GetCategories)
			categories.GET("/tree", controller.GetCategoryTree)
			categories.GET("/:id", controller.GetCategory)
			categories.GET("/:id/children", controller.GetSubCategories)

			// 管理员路由
			adminCategories := categories.Group("")
			adminCategories.Use(middleware.RequireAdmin())
			{
				adminCategories.POST("", controller.CreateCategory)
				adminCategories.PUT("/:id", controller.UpdateCategory)
				adminCategories.DELETE("/:id", controller.DeleteCategory)
				adminCategories.PUT("/:id/status", controller.UpdateCategoryStatus)
			}
		}

		// 购物车路由 - 待实现
		/*
			cart := api.Group("/cart")
			cart.Use(middleware.JWTAuth())
			{
				cart.GET("", controller.GetCartList)
				cart.POST("", controller.AddToCart)
				cart.PUT("/:id/quantity", controller.UpdateCartQuantity)
				cart.PUT("/:id/selected", controller.UpdateCartSelected)
				cart.DELETE("/:id", controller.DeleteCartItem)
				cart.DELETE("/selected", controller.DeleteSelectedItems)
				cart.POST("/checkout", controller.Checkout)
			}
		*/

		// 订单路由 - 待实现
		/*
			orders := api.Group("/orders")
			orders.Use(middleware.JWTAuth())
			{
				orders.GET("", controller.GetOrderList)
				orders.GET("/:id", controller.GetOrderDetail)
				orders.POST("", controller.CreateOrder)
				orders.PUT("/:id/cancel", controller.CancelOrder)
				orders.PUT("/:id/receive", controller.ReceiveOrder)
				orders.DELETE("/:id", controller.DeleteOrder)

				// 管理员路由
				adminOrders := orders.Group("")
				adminOrders.Use(middleware.RequireAdmin())
				{
					adminOrders.PUT("/:id/status", controller.UpdateOrderStatus)
					adminOrders.PUT("/:id/ship", controller.ShipOrder)
					adminOrders.GET("/statistics", controller.GetOrderStatistics)
				}
			}
		*/

		// 优惠券路由 - 待实现
		/*
			coupons := api.Group("/coupons")
			{
				coupons.GET("", controller.GetCouponList)

				// 用户优惠券路由
				userCoupons := api.Group("/my")
				userCoupons.Use(middleware.JWTAuth())
				{
					userCoupons.GET("", controller.GetUserCoupons)
					userCoupons.POST("/:id/receive", controller.ReceiveCoupon)
					userCoupons.GET("/:id/validate", controller.ValidateCoupon)
				}

				// 管理员路由
				adminCoupons := coupons.Group("")
				adminCoupons.Use(middleware.RequireAdmin())
				{
					adminCoupons.POST("", controller.CreateCoupon)
					adminCoupons.PUT("/:id", controller.UpdateCoupon)
					adminCoupons.DELETE("/:id", controller.DeleteCoupon)
					adminCoupons.PUT("/:id/status", controller.UpdateCouponStatus)
				}
			}
		*/

		// 搜索路由 - 待实现
		// api.GET("/search", controller.SearchProducts)

		// 上传路由 - 待实现
		// api.POST("/upload", middleware.JWTAuth(), controller.UploadFile)
	}

	// 静态文件服务
	r.Static("/static", "./static")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")

	return r
}

// healthCheck 健康检查
func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"service": "online-mall",
		"time":    "2024",
	})
}
