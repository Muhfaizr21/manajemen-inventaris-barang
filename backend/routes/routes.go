package routes

import (
	"backend/controllers"
	"backend/middlewares"
	"backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter mengonfigurasi seluruh routing HTTP dengan injeksi Controller dan Middleware
func SetupRouter(
	authController *controllers.AuthController,
	itemController *controllers.ItemController,
	authService services.AuthService,
) *gin.Engine {
	r := gin.Default()

	// Health check & informasi layanan
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "Sistem Manajemen Inventaris Barang API",
			"version": "1.0.0",
			"status":  "healthy",
		})
	})

	// Endpoint Authentication (Root level sesuai spesifikasi soal & /api prefix untuk fleksibilitas)
	r.POST("/register", authController.Register)
	r.POST("/login", authController.Login)

	// Endpoint CRUD Barang (Wajib Login / JWT Protected)
	items := r.Group("/items")
	items.Use(middlewares.AuthMiddleware(authService))
	{
		items.POST("", itemController.CreateItem)
		items.GET("", itemController.GetItems)
		items.GET("/:id", itemController.GetItemByID)
		items.PUT("/:id", itemController.UpdateItem)
		items.DELETE("/:id", itemController.DeleteItem)
	}

	// Alias dengan prefix /api untuk kompatibilitas standar RESTful API
	api := r.Group("/api")
	{
		api.POST("/register", authController.Register)
		api.POST("/login", authController.Login)

		apiItems := api.Group("/items")
		apiItems.Use(middlewares.AuthMiddleware(authService))
		{
			apiItems.POST("", itemController.CreateItem)
			apiItems.GET("", itemController.GetItems)
			apiItems.GET("/:id", itemController.GetItemByID)
			apiItems.PUT("/:id", itemController.UpdateItem)
			apiItems.DELETE("/:id", itemController.DeleteItem)
		}
	}

	return r
}
