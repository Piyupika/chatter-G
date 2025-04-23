package main

import (
	"log"
	"os"

	"github.com/Piyu-Pika/godzilla-go/internal/database"
	"github.com/Piyu-Pika/godzilla-go/internal/routes"
	"github.com/joho/godotenv"

	// "github.com/Piyu-Pika/godzilla-go/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// Set the Gin mode based on the environment variable GIN_MODE

	if os.Getenv("GIN_MODE") == "debug" {

		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	// Create default gin engine
	r := gin.Default()

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Set the log output to a file
	if os.Getenv("GIN_MODE") == "debug" {
		logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("Failed to open log file:", err)
		}
		defer logFile.Close()

		log.SetOutput(logFile)
	}

	// Setup routes
	routes.SetupRoutes(r)

	// Initialize the database connection
	database.Init()

	// Check if the database connection is successful
	if database.DB == nil {
		log.Fatal("Failed to initialize database connection")
	}
	log.Println("Database connection initialized successfully")

	//Database Close

	defer database.CloseDB()

	if os.Getenv("GIN_MODE") == "debug" {
		// ws := services.NewWebSocketManager()
		// go ws.Run()
	}

	// Run the server on port 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
