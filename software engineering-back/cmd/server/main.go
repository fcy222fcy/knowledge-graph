package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"software_engineering/internal/database"
	"software_engineering/internal/middleware"
	"software_engineering/internal/api"
	"software_engineering/internal/repository/seed"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found, using environment variables")
	}

	// Connect to database
	database.Connect()

	// AutoMigrate tables
	database.AutoMigrate()

	// Seed demo data
	seed.SeedAll()

	// Setup Gin router
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// Register all routes
	api.SetupRoutes(r)

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
