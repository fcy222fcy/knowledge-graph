package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	Neo4jURI   string
	Neo4jUser  string
	Neo4jPass  string
	JWTSecret  string
}

var AppConfig Config

func Load() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	AppConfig = Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "software_engineering"),
		Neo4jURI:   getEnv("NEO4J_URI", ""),
		Neo4jUser:  getEnv("NEO4J_USER", "neo4j"),
		Neo4jPass:  getEnv("NEO4J_PASSWORD", ""),
		JWTSecret:  getEnv("JWT_SECRET", "software-engineering-qa-platform-secret-key"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
