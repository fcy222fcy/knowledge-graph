package config

import (
	"os"
	"testing"
)

func TestLoad_DefaultValues(t *testing.T) {
	// 清除所有相关环境变量
	envVars := []string{
		"SERVER_PORT", "DB_HOST", "DB_PORT", "DB_USER",
		"DB_PASSWORD", "DB_NAME", "NEO4J_URI", "NEO4J_USER",
		"NEO4J_PASSWORD", "JWT_SECRET",
	}
	for _, v := range envVars {
		os.Unsetenv(v)
	}

	Load()

	// 验证默认值
	if AppConfig.ServerPort != "8080" {
		t.Errorf("ServerPort = %v, want '8080'", AppConfig.ServerPort)
	}
	if AppConfig.DBHost != "localhost" {
		t.Errorf("DBHost = %v, want 'localhost'", AppConfig.DBHost)
	}
	if AppConfig.DBPort != "3306" {
		t.Errorf("DBPort = %v, want '3306'", AppConfig.DBPort)
	}
	if AppConfig.DBUser != "root" {
		t.Errorf("DBUser = %v, want 'root'", AppConfig.DBUser)
	}
	if AppConfig.DBName != "software_engineering" {
		t.Errorf("DBName = %v, want 'software_engineering'", AppConfig.DBName)
	}
	if AppConfig.Neo4jUser != "neo4j" {
		t.Errorf("Neo4jUser = %v, want 'neo4j'", AppConfig.Neo4jUser)
	}
}

func TestLoad_EnvironmentVariables(t *testing.T) {
	// 设置测试环境变量
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "3307")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("NEO4j_URI", "bolt://neo4j.example.com:7687")
	os.Setenv("NEO4J_USER", "neo4j_admin")
	os.Setenv("NEO4J_PASSWORD", "neo4j_pass")
	os.Setenv("JWT_SECRET", "my-secret-key")

	defer func() {
		// 清理环境变量
		for _, v := range []string{
			"SERVER_PORT", "DB_HOST", "DB_PORT", "DB_USER",
			"DB_PASSWORD", "DB_NAME", "NEO4j_URI", "NEO4J_USER",
			"NEO4J_PASSWORD", "JWT_SECRET",
		} {
			os.Unsetenv(v)
		}
	}()

	Load()

	// 验证环境变量值
	if AppConfig.ServerPort != "9090" {
		t.Errorf("ServerPort = %v, want '9090'", AppConfig.ServerPort)
	}
	if AppConfig.DBHost != "db.example.com" {
		t.Errorf("DBHost = %v, want 'db.example.com'", AppConfig.DBHost)
	}
	if AppConfig.DBPort != "3307" {
		t.Errorf("DBPort = %v, want '3307'", AppConfig.DBPort)
	}
	if AppConfig.DBUser != "testuser" {
		t.Errorf("DBUser = %v, want 'testuser'", AppConfig.DBUser)
	}
	if AppConfig.DBPassword != "testpass" {
		t.Errorf("DBPassword = %v, want 'testpass'", AppConfig.DBPassword)
	}
	if AppConfig.DBName != "testdb" {
		t.Errorf("DBName = %v, want 'testdb'", AppConfig.DBName)
	}
	if AppConfig.Neo4jURI != "bolt://neo4j.example.com:7687" {
		t.Errorf("Neo4jURI = %v, want 'bolt://neo4j.example.com:7687'", AppConfig.Neo4jURI)
	}
	if AppConfig.Neo4jUser != "neo4j_admin" {
		t.Errorf("Neo4jUser = %v, want 'neo4j_admin'", AppConfig.Neo4jUser)
	}
	if AppConfig.Neo4jPass != "neo4j_pass" {
		t.Errorf("Neo4jPass = %v, want 'neo4j_pass'", AppConfig.Neo4jPass)
	}
	if AppConfig.JWTSecret != "my-secret-key" {
		t.Errorf("JWTSecret = %v, want 'my-secret-key'", AppConfig.JWTSecret)
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		value        string
		defaultValue string
		want         string
	}{
		{
			name:         "环境变量存在",
			key:          "TEST_GETENV_EXISTS",
			value:        "custom-value",
			defaultValue: "default-value",
			want:         "custom-value",
		},
		{
			name:         "环境变量不存在，使用默认值",
			key:          "TEST_GETENV_NOT_EXISTS",
			value:        "",
			defaultValue: "default-value",
			want:         "default-value",
		},
		{
			name:         "空默认值",
			key:          "TEST_GETENV_EMPTY_DEFAULT",
			value:        "",
			defaultValue: "",
			want:         "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			got := getEnv(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("getEnv(%v, %v) = %v, want %v", tt.key, tt.defaultValue, got, tt.want)
			}
		})
	}
}
