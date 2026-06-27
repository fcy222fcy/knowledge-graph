package admin

import (
	"github.com/gin-gonic/gin"
	"software_engineering/pkg/config"
	"software_engineering/pkg/response"
)

// SystemConfig 系统配置响应
type SystemConfig struct {
	OllamaURL   string `json:"ollama_url"`
	OllamaModel string `json:"ollama_model"`
	ServerPort  string `json:"server_port"`
}

// GetSystemConfig 获取系统配置
func GetSystemConfig(c *gin.Context) {
	cfg := config.AppConfig

	response.Success(c, SystemConfig{
		OllamaURL:   cfg.OllamaURL,
		OllamaModel: cfg.OllamaModel,
		ServerPort:  cfg.ServerPort,
	})
}
