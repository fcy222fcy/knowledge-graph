# Core Skeleton Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a compilable Go backend skeleton with all GORM models, database connection, CORS middleware, route registration with placeholder handlers, seed data, and a working health-check endpoint.

**Architecture:** Follow the MVC-ish layered architecture defined in `AGENTS.md`. The skeleton sets up `cmd/server` entry point, `internal/database` for MySQL connection and AutoMigrate, `internal/model` for all GORM models, `internal/middleware` for CORS, `internal/routes` for route registration, and `internal/seed` for demo data. No business logic — only structural wiring.

**Tech Stack:** Go 1.24, Gin v1.10.1, GORM v1.25.12, MySQL (via gorm.io/driver/mysql), godotenv for .env loading

---

## File Structure

| Path | Responsibility |
|------|---------------|
| `cmd/server/main.go` | Entry point: load env, connect DB, migrate, seed, start server |
| `internal/database/database.go` | MySQL connection setup via GORM |
| `internal/database/migrate.go` | AutoMigrate all models |
| `internal/model/model.go` | All GORM data models (User, Document, KnowledgePoint, KnowledgeRelation, KnowledgeBuild, Question, Quiz, AskSession, AskMessage) |
| `internal/middleware/cors.go` | CORS middleware for Gin |
| `internal/middleware/auth.go` | JWT auth middleware stub (always passes) |
| `internal/routes/routes.go` | Register all API route groups with placeholder handlers |
| `internal/seed/seed.go` | Seed demo users, knowledge points, questions |

---

## Task 1: Add godotenv dependency

**Files:**
- Modify: `go.mod`

- [ ] **Step 1: Install godotenv**

```bash
cd "E:/goCode/goFile/software engineering/software engineering-back"
go get github.com/joho/godotenv
```

Expected: `go.mod` and `go.sum` updated with godotenv entry.

- [ ] **Step 2: Verify dependency resolves**

```bash
go mod tidy
```

Expected: No errors.

- [ ] **Step 3: Commit**

```bash
git add go.mod go.sum
git commit -m "chore: add godotenv dependency for .env loading"
```

---

## Task 2: GORM data models

**Files:**
- Create: `internal/model/model.go`

- [ ] **Step 1: Create model package with all data models**

Create `internal/model/model.go` with the following content. These models match the API spec in `doc/api.md` and the architecture in `AGENTS.md`.

```go
package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel embeds gorm.Model for common fields.
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ─── User ─────────────────────────────────────────────

type User struct {
	BaseModel
	Username string `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password string `gorm:"size:255;not null" json:"-"`
	Email    string `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Nickname string `gorm:"size:50" json:"nickname"`
	Avatar   string `gorm:"size:255" json:"avatar"`
	Status   int    `gorm:"default:1" json:"status"` // 1=active, 0=disabled
}

// ─── Document ─────────────────────────────────────────

type Document struct {
	BaseModel
	Title       string `gorm:"size:200;not null" json:"title"`
	Description string `gorm:"size:500" json:"description"`
	Filename    string `gorm:"size:200;not null" json:"filename"`
	FilePath    string `gorm:"size:500;not null" json:"-"`
	FileSize    int64  `json:"file_size"`
	FileType    string `gorm:"size:20" json:"file_type"`
	Content     string `gorm:"type:longtext" json:"-"`
	Status      string `gorm:"size:20;default:pending" json:"status"` // pending/processing/completed/failed
}

// ─── Knowledge ────────────────────────────────────────

type KnowledgePoint struct {
	BaseModel
	Name        string `gorm:"size:100;not null" json:"name"`
	Description string `gorm:"size:500" json:"description"`
	DocumentID  uint   `json:"document_id"`
	Category    string `gorm:"size:50" json:"category"`
}

type KnowledgeRelation struct {
	BaseModel
	SourceID     uint   `json:"source_id"`
	TargetID     uint   `json:"target_id"`
	RelationType string `gorm:"size:20;not null" json:"relation_type"` // RELATED/DEPENDS_ON/PART_OF
	Description  string `gorm:"size:500" json:"description"`
}

type KnowledgeBuild struct {
	BaseModel
	DocumentIDs    string `gorm:"size:500" json:"document_ids"` // comma-separated
	CreatedPoints  int    `json:"created_points"`
	CreatedRelations int  `json:"created_relations"`
	ChunkCount     int    `json:"chunk_count"`
	VectorCount    int    `json:"vector_count"`
	Status         string `gorm:"size:20;default:completed" json:"status"`
	Message        string `gorm:"size:500" json:"message"`
}

// ─── Question ─────────────────────────────────────────

type Question struct {
	BaseModel
	Title             string `gorm:"size:500;not null" json:"title"`
	Type              string `gorm:"size:20;not null" json:"type"` // single/multiple
	Difficulty        string `gorm:"size:20;not null" json:"difficulty"` // easy/medium/hard
	KnowledgePointID  uint   `json:"knowledge_point_id"`
	Options           string `gorm:"type:text" json:"-"` // JSON array stored as text
	Answer            string `gorm:"size:20;not null" json:"answer"`
	Explanation       string `gorm:"type:text" json:"explanation"`
}

// ─── Quiz ─────────────────────────────────────────────

type Quiz struct {
	BaseModel
	QuestionID uint   `json:"question_id"`
	UserID     uint   `json:"user_id"`
	UserAnswer string `gorm:"size:20;not null" json:"user_answer"`
	IsCorrect  bool   `json:"is_correct"`
}

// ─── Ask (Q&A) ────────────────────────────────────────

type AskSession struct {
	BaseModel
	UserID  uint   `json:"user_id"`
	Title   string `gorm:"size:200" json:"title"`
}

type AskMessage struct {
	BaseModel
	SessionID   uint   `json:"session_id"`
	Role        string `gorm:"size:20;not null" json:"role"` // user/assistant
	Content     string `gorm:"type:text" json:"content"`
	Confidence  float64 `json:"confidence"`
}
```

- [ ] **Step 2: Verify compilation**

```bash
go build ./internal/model/...
```

Expected: No errors.

- [ ] **Step 3: Commit**

```bash
git add internal/model/model.go
git commit -m "feat: add all GORM data models"
```

---

## Task 3: Database connection and AutoMigrate

**Files:**
- Create: `internal/database/database.go`
- Create: `internal/database/migrate.go`

- [ ] **Step 1: Create database connection package**

Create `internal/database/database.go`:

```go
package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to get underlying sql.DB: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("database connected successfully")
}
```

- [ ] **Step 2: Create AutoMigrate function**

Create `internal/database/migrate.go`:

```go
package database

import (
	"log"
	"software_engineering/internal/model"
)

func AutoMigrate() {
	err := DB.AutoMigrate(
		&model.User{},
		&model.Document{},
		&model.KnowledgePoint{},
		&model.KnowledgeRelation{},
		&model.KnowledgeBuild{},
		&model.Question{},
		&model.Quiz{},
		&model.AskSession{},
		&model.AskMessage{},
	)
	if err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}
	log.Println("database migration completed")
}
```

- [ ] **Step 3: Verify compilation**

```bash
go build ./internal/database/...
```

Expected: No errors.

- [ ] **Step 4: Commit**

```bash
git add internal/database/
git commit -m "feat: add database connection and AutoMigrate"
```

---

## Task 4: CORS middleware

**Files:**
- Create: `internal/middleware/cors.go`

- [ ] **Step 1: Create CORS middleware**

Create `internal/middleware/cors.go`:

```go
package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
```

- [ ] **Step 2: Verify compilation**

```bash
go build ./internal/middleware/...
```

Expected: No errors.

- [ ] **Step 3: Commit**

```bash
git add internal/middleware/cors.go
git commit -m "feat: add CORS middleware"
```

---

## Task 5: Auth middleware stub

**Files:**
- Create: `internal/middleware/auth.go`

- [ ] **Step 1: Create auth middleware stub**

Create `internal/middleware/auth.go`. This is a stub that always passes — real JWT logic will be added later.

```go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a stub that currently allows all requests.
// TODO: implement JWT token validation.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Stub: skip auth check for now
		c.Next()
	}
}

// respondError writes a JSON error response and aborts.
func respondError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
		"data":    nil,
	})
	c.Abort()
}

// RequireAuth is a placeholder that checks for token presence (stub).
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			respondError(c, http.StatusUnauthorized, "未授权")
			return
		}
		// Stub: accept any non-empty token
		c.Set("user_id", uint(1))
		c.Next()
	}
}
```

- [ ] **Step 2: Verify compilation**

```bash
go build ./internal/middleware/...
```

Expected: No errors.

- [ ] **Step 3: Commit**

```bash
git add internal/middleware/auth.go
git commit -m "feat: add auth middleware stub"
```

---

## Task 6: Routes with placeholder handlers

**Files:**
- Create: `internal/routes/routes.go`

- [ ] **Step 1: Create route registration with all placeholder endpoints**

Create `internal/routes/routes.go`. This registers all API groups from `doc/api.md` with stub handlers that return the standard response format.

```go
package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/middleware"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		// Health check
		api.GET("/health", healthCheck)

		// Auth (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", stubHandler("用户注册"))
			auth.POST("/login", stubHandler("用户登录"))
			auth.POST("/refresh", stubHandler("刷新token"))
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.RequireAuth())
		{
			// Users
			users := protected.Group("/users")
			{
				users.GET("/profile", stubHandler("获取当前用户信息"))
				users.PUT("/profile", stubHandler("更新用户信息"))
				users.POST("/password", stubHandler("修改密码"))
				users.GET("/list", paginatedStub("用户列表"))
			}

			// Documents
			docs := protected.Group("/documents")
			{
				docs.POST("", stubHandler("上传文档"))
				docs.GET("", paginatedStub("文档列表"))
				docs.GET("/:id", stubHandler("文档详情"))
				docs.PUT("/:id", stubHandler("更新文档信息"))
				docs.DELETE("/:id", stubHandler("删除文档"))
				docs.GET("/:id/content", stubHandler("获取文档内容"))
			}

			// Knowledge points
			kp := protected.Group("/knowledge")
			{
				kp.GET("/points", paginatedStub("知识点列表"))
				kp.GET("/points/:id", stubHandler("知识点详情"))
				kp.POST("/points", stubHandler("新增知识点"))
				kp.PUT("/points/:id", stubHandler("更新知识点"))
				kp.DELETE("/points/:id", stubHandler("删除知识点"))
				kp.GET("/relations", paginatedStub("关系列表"))
				kp.POST("/relations", stubHandler("新增关系"))
				kp.PUT("/relations/:id", stubHandler("更新关系"))
				kp.DELETE("/relations/:id", stubHandler("删除关系"))
			}

			// Graph
			graph := protected.Group("/graph")
			{
				graph.GET("", stubHandler("获取图谱数据"))
				graph.POST("/build", stubHandler("从文档构建图谱"))
				graph.GET("/build/latest", stubHandler("最近构建结果"))
				graph.GET("/build/history", paginatedStub("构建历史记录"))
			}

			// Questions
			q := protected.Group("/questions")
			{
				q.GET("", paginatedStub("题目列表"))
				q.GET("/:id", stubHandler("题目详情"))
				q.POST("", stubHandler("新增题目"))
				q.PUT("/:id", stubHandler("更新题目"))
				q.DELETE("/:id", stubHandler("删除题目"))
			}

			// Quizzes
			quiz := protected.Group("/quizzes")
			{
				quiz.POST("/submit", stubHandler("提交答题"))
				quiz.GET("/history", paginatedStub("答题历史"))
				quiz.GET("/:id", stubHandler("答题详情"))
			}

			// Ask (Q&A)
			ask := protected.Group("/ask")
			{
				ask.POST("/sessions", stubHandler("新建问答会话"))
				ask.GET("/sessions", paginatedStub("会话列表"))
				ask.GET("/sessions/:id/messages", paginatedStub("会话消息列表"))
				ask.POST("", stubHandler("提问"))
				ask.GET("/history", paginatedStub("问答历史"))
			}

			// Analytics
			analytics := protected.Group("/analytics")
			{
				analytics.GET("/overview", stubHandler("总览统计"))
				analytics.GET("/hot-knowledge-points", stubHandler("热门知识点"))
				analytics.GET("/knowledge-mastery", stubHandler("知识点掌握度"))
				analytics.GET("/weak-points", stubHandler("薄弱知识点"))
				analytics.GET("/trends", stubHandler("趋势数据"))
			}
		}
	}
}

// healthCheck returns service status.
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"status":  "ok",
			"service": "software-engineering-backend",
		},
	})
}

// stubHandler returns a placeholder handler for unimplemented endpoints.
func stubHandler(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": name + " - 待实现",
			"data":    nil,
		})
	}
}

// paginatedStub returns a placeholder handler that returns an empty paginated list.
func paginatedStub(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": name + " - 待实现",
			"data": gin.H{
				"list":       []interface{}{},
				"total":      0,
				"page":       page,
				"size":       size,
				"total_page": 0,
			},
		})
	}
}
```

- [ ] **Step 2: Verify compilation**

```bash
go build ./internal/routes/...
```

Expected: No errors.

- [ ] **Step 3: Commit**

```bash
git add internal/routes/routes.go
git commit -m "feat: add route registration with placeholder handlers"
```

---

## Task 7: Seed data

**Files:**
- Create: `internal/seed/seed.go`

- [ ] **Step 1: Create seed data package**

Create `internal/seed/seed.go`. This seeds demo users, knowledge points, and questions when tables are empty.

```go
package seed

import (
	"log"
	"software_engineering/internal/database"
	"software_engineering/internal/model"
)

func SeedAll() {
	seedUsers()
	seedKnowledgePoints()
	seedQuestions()
}

func seedUsers() {
	var count int64
	database.DB.Model(&model.User{}).Count(&count)
	if count > 0 {
		return
	}

	users := []model.User{
		{Username: "student001", Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", Email: "student001@example.com", Nickname: "张三", Status: 1},
		{Username: "student002", Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", Email: "student002@example.com", Nickname: "李四", Status: 1},
		{Username: "teacher001", Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", Email: "teacher001@example.com", Nickname: "王老师", Status: 1},
	}

	if err := database.DB.Create(&users).Error; err != nil {
		log.Printf("seed users failed: %v", err)
		return
	}
	log.Println("seeded 3 demo users")
}

func seedKnowledgePoints() {
	var count int64
	database.DB.Model(&model.KnowledgePoint{}).Count(&count)
	if count > 0 {
		return
	}

	points := []model.KnowledgePoint{
		{Name: "需求分析", Description: "识别和确认用户需求的过程", Category: "需求相关"},
		{Name: "软件测试", Description: "验证软件是否满足需求的过程", Category: "测试相关"},
		{Name: "软件生命周期", Description: "软件从提出到废弃的整个过程", Category: "基础概念"},
		{Name: "编码实现", Description: "将设计转化为可执行代码的过程", Category: "开发相关"},
		{Name: "项目管理", Description: "对软件项目进行计划、组织和控制", Category: "管理相关"},
	}

	if err := database.DB.Create(&points).Error; err != nil {
		log.Printf("seed knowledge points failed: %v", err)
		return
	}

	// Seed relations
	relations := []model.KnowledgeRelation{
		{SourceID: 1, TargetID: 2, RelationType: "DEPENDS_ON", Description: "需求分析是软件测试的前置环节"},
		{SourceID: 1, TargetID: 4, RelationType: "DEPENDS_ON", Description: "需求分析完成后进入编码实现"},
		{SourceID: 3, TargetID: 1, RelationType: "PART_OF", Description: "需求分析是软件生命周期的一个阶段"},
		{SourceID: 5, TargetID: 3, RelationType: "RELATED", Description: "项目管理贯穿整个软件生命周期"},
	}

	if err := database.DB.Create(&relations).Error; err != nil {
		log.Printf("seed knowledge relations failed: %v", err)
		return
	}
	log.Println("seeded 5 knowledge points and 4 relations")
}

func seedQuestions() {
	var count int64
	database.DB.Model(&model.Question{}).Count(&count)
	if count > 0 {
		return
	}

	questions := []model.Question{
		{
			Title:            "以下哪个不是需求分析的活动？",
			Type:             "single",
			Difficulty:       "easy",
			KnowledgePointID: 1,
			Options:          `[{"key":"A","value":"需求获取"},{"key":"B","value":"需求分析"},{"key":"C","value":"代码编写"},{"key":"D","value":"需求验证"}]`,
			Answer:           "C",
			Explanation:      "代码编写属于编码阶段，不是需求分析的活动",
		},
		{
			Title:            "黑盒测试的主要关注点是什么？",
			Type:             "single",
			Difficulty:       "easy",
			KnowledgePointID: 2,
			Options:          `[{"key":"A","value":"程序内部逻辑"},{"key":"B","value":"程序外部功能"},{"key":"C","value":"代码覆盖率"},{"key":"D","value":"算法效率"}]`,
			Answer:           "B",
			Explanation:      "黑盒测试关注程序的外部功能，不关心内部实现细节",
		},
		{
			Title:            "软件生命周期的第一个阶段是？",
			Type:             "single",
			Difficulty:       "easy",
			KnowledgePointID: 3,
			Options:          `[{"key":"A","value":"编码"},{"key":"B","value":"测试"},{"key":"C","value":"需求分析"},{"key":"D","value":"维护"}]`,
			Answer:           "C",
			Explanation:      "软件生命周期的第一个阶段是需求分析",
		},
	}

	if err := database.DB.Create(&questions).Error; err != nil {
		log.Printf("seed questions failed: %v", err)
		return
	}
	log.Println("seeded 3 demo questions")
}
```

- [ ] **Step 2: Verify compilation**

```bash
go build ./internal/seed/...
```

Expected: No errors.

- [ ] **Step 3: Commit**

```bash
git add internal/seed/seed.go
git commit -m "feat: add seed data for users, knowledge points, and questions"
```

---

## Task 8: Entry point

**Files:**
- Create: `cmd/server/main.go`

- [ ] **Step 1: Create main.go entry point**

Create `cmd/server/main.go`:

```go
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"software_engineering/internal/database"
	"software_engineering/internal/routes"
	"software_engineering/internal/seed"
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

	// Register all routes
	routes.SetupRoutes(r)

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
```

- [ ] **Step 2: Verify full project compilation**

```bash
go build ./cmd/server/...
```

Expected: No errors. Binary produced.

- [ ] **Step 3: Verify with go vet**

```bash
go vet ./...
```

Expected: No warnings or errors.

- [ ] **Step 4: Commit**

```bash
git add cmd/server/main.go
git commit -m "feat: add server entry point with env loading, DB, migration, and seeding"
```

---

## Task 9: Full build verification

**Files:** None (verification only)

- [ ] **Step 1: Tidy modules**

```bash
go mod tidy
```

Expected: No errors.

- [ ] **Step 2: Build entire project**

```bash
go build ./...
```

Expected: No errors.

- [ ] **Step 3: Run vet on all packages**

```bash
go vet ./...
```

Expected: No warnings.

- [ ] **Step 4: Final commit if any cleanup was needed**

```bash
git add -A
git commit -m "chore: verify full project builds cleanly"
```

(Only run this commit if `go mod tidy` or other steps produced changes.)
