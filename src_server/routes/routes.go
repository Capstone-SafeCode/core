package routes

import (
	"net/http"
	"test_capstone/src_server/controllers"
	middlewares "test_capstone/src_server/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Middleware CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	r.POST("/register", controllers.CreateUser)
	r.POST("/login", controllers.LogUser)

	// Routes protégées
	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/users", controllers.GetUsers)
		protected.GET("/me", controllers.GetMe)

		protected.POST("/upload", controllers.UploadFile)
		protected.POST("/analyse", controllers.Analyse)
	}

	return r
}
