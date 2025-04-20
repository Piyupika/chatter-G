package database

import (
	"context"
	"log"

	"github.com/Piyu-Pika/godzilla-go/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var MongoClient *mongo.Client
var ChatCollection *mongo.Collection

func InitDB() {
	dsn := `mkdir -p $env:appdata\postgresql\; Invoke-WebRequest -Uri https://cockroachlabs.cloud/clusters/4912e37b-f4b6-4f64-a957-9586696914b8/cert -OutFile $env:appdata\postgresql\root.crt` // Replace with your CockroachDB DSN
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to CockroachDB:", err)
	}
	// Auto migrate models
	DB.AutoMigrate(&models.User{})
}

func InitMongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://piyush:<db_password>@cluster0.gjp6pft.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0") // Replace with your MongoDB URI
	var err error
	MongoClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	err = MongoClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}
	ChatCollection = MongoClient.Database("chatdb").Collection("messages")
}

func CloseMongoDB() {
	MongoClient.Disconnect(context.TODO())
}

func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to close CockroachDB connection:", err)
	}
	sqlDB.Close()
	CloseMongoDB()
	log.Println("Database connections closed successfully")
}

func Init() {
	InitDB()
	InitMongoDB()
	log.Println("Database connections initialized successfully")
}
