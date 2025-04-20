package routes

import (
	"github.com/Piyu-Pika/godzilla-go/internal/handlers"
	"github.com/Piyu-Pika/godzilla-go/internal/services"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// Authentication routes
	r.GET("/health", handlers.HealthCheck)
	r.POST("/save-data", handlers.SaveData)
	// User listing (protected route)
	r.GET("/users", handlers.GetUsers)

	// WebSocket route
	r.GET("/ws", func(c *gin.Context) {
		ws := services.NewWebSocketManager() // In practice, use a singleton or dependency injection
		ws.HandleWebSocketConnections(c.Writer, c.Request)
	})
}
