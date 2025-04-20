package main

import (
	"log"
	"os"

	"github.com/Piyu-Pika/godzilla-go/internal/routes"
	"github.com/Piyu-Pika/godzilla-go/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create default gin engine
	r := gin.Default()

	// Set the log output to a file
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	// Setup routes
	routes.SetupRoutes(r)

	ws := services.NewWebSocketManager()
	go ws.Run()

	//graceful shutdown

	// Run the server on port 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
