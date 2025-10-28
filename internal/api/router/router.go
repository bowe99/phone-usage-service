package router

import (
	"github.com/bowe99/phone-usage-service/internal/api/handler"
	"github.com/bowe99/phone-usage-service/internal/infra/database"
	"github.com/gin-gonic/gin"
)

func SetupRouter(db *database.MongoDB, ginMode string, userHandler *handler.UserHandler) *gin.Engine {
	gin.SetMode(ginMode)
	router := gin.New()

	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		if err := db.HealthCheck(c.Request.Context()); err != nil {
			c.JSON(500, gin.H{"status": "unhealth", "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "heatlhy"})
	})

	users := router.Group("/api/users")
	{
		users.POST("", userHandler.CreateUser)
		users.PUT("/:id", userHandler.UpdateUserProfile)
	}

	return router
}