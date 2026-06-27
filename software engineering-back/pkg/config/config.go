package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config 应用配置结构体，包含服务器端口、数据库连接、Neo4j、JWT 和 AI 相关配置
type Config struct {
	ServerPort string // 服务器端口
	DBHost     string // 数据库主机
	DBPort     string // 数据库端口
	DBUser     string // 数据库用户名
	DBPassword string // 数据库密码
	DBName     string // 数据库名称
	Neo4jURI   string // Neo4j 连接URI
	Neo4jUser  string // Neo4j 用户名
	Neo4jPass  string // Neo4j 密码
	JWTSecret  string // JWT 签名密钥

	// AI 配置 (Ollama)
	OllamaURL            string // Ollama 服务地址
	OllamaModel          string // Ollama 生成模型
	OllamaEmbeddingModel string // Ollama 嵌入模型
}

// AppConfig 全局配置实例，程序启动时通过 Load() 初始化
var AppConfig Config

// Load 从环境变量加载配置，优先读取 .env 文件
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

		// AI 配置 (Ollama)
		OllamaURL:            getEnv("OLLAMA_URL", "http://localhost:11434"),
		OllamaModel:          getEnv("OLLAMA_MODEL", "qwen3:8b"),
		OllamaEmbeddingModel: getEnv("OLLAMA_EMBEDDING_MODEL", "nomic-embed-text"),
	}
}

// getEnv 获取环境变量值，不存在时返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
