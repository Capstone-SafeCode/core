package routes

import (
	"net/http"
	"test_capstone/src_server/controllers"

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
	r.GET("/users", controllers.GetUsers)

	r.POST("/upload", controllers.UploadFile)
	r.POST("/analyse", controllers.StartAnalyse)
	return r
}
