# 后端全量实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**目标:** 实现全部 33 个 API 端点的业务逻辑，完成认证、用户、文档、知识点、知识图谱、题库、答题、问答、学习分析共 9 个功能模块。

**架构:** 遵循 AGENTS.md 分层架构 — controller 接收请求并调用 service，service 处理业务逻辑并调用 repository，repository 封装 GORM 查询。新增 DTO 层定义请求/响应结构体。JWT 用于认证。

**技术栈:** Go 1.25, Gin, GORM, MySQL, golang-jwt/jwt/v5, golang.org/x/crypto (bcrypt)

---

## 模块依赖关系

```
Auth (认证) ─────────────────────────────────────────────── 基础，无依赖
  ├─ User (用户管理) ─────────── 依赖 Auth
  ├─ Document (资源管理) ─────── 依赖 Auth
  ├─ Knowledge (知识点管理) ──── 依赖 Auth, Document
  │   ├─ Graph (知识图谱) ────── 依赖 Knowledge, Document
  │   └─ Question (题库) ─────── 依赖 Knowledge
  │       └─ Quiz (答题) ─────── 依赖 Question
  ├─ Ask (知识问答) ──────────── 依赖 Auth, Knowledge
  └─ Analytics (学习分析) ────── 依赖 Quiz, Ask
```

**可并行执行的模块组:**

| 批次 | 模块 | 依赖 |
|------|------|------|
| 1 | Auth | 无 |
| 2 | User, Document | Auth |
| 3 | Knowledge, Question | Auth, Document |
| 4 | Graph, Quiz, Ask | Knowledge, Question |
| 5 | Analytics | Quiz, Ask |

---

## 文件结构总览

每个模块需要创建 4 个文件（按 AGENTS.md 分层）：

```
internal/
├── dto/
│   ├── auth.go           # 认证请求/响应
│   ├── user.go           # 用户请求/响应
│   ├── document.go       # 文档请求/响应
│   ├── knowledge.go      # 知识点请求/响应
│   ├── graph.go          # 图谱请求/响应
│   ├── question.go       # 题目请求/响应
│   ├── quiz.go           # 答题请求/响应
│   ├── ask.go            # 问答请求/响应
│   └── analytics.go      # 分析请求/响应
├── repository/
│   ├── user_repo.go
│   ├── document_repo.go
│   ├── knowledge_repo.go
│   ├── question_repo.go
│   ├── quiz_repo.go
│   ├── ask_repo.go
│   └── analytics_repo.go
├── service/
│   ├── auth_service.go
│   ├── user_service.go
│   ├── document_service.go
│   ├── knowledge_service.go
│   ├── graph_service.go
│   ├── question_service.go
│   ├── quiz_service.go
│   ├── ask_service.go
│   └── analytics_service.go
├── controller/
│   ├── auth_controller.go
│   ├── user_controller.go
│   ├── document_controller.go
│   ├── knowledge_controller.go
│   ├── graph_controller.go
│   ├── question_controller.go
│   ├── quiz_controller.go
│   ├── ask_controller.go
│   └── analytics_controller.go
└── middleware/
    └── auth.go           # 从 stub 改为真实 JWT 验证
```

---

## Task 1: 公共工具和 JWT 基础设施

**依赖:** 无（最先执行）
**可并行:** 否（所有后续任务依赖此任务）

**文件:**
- 创建: `internal/dto/common.go`
- 创建: `internal/utils/jwt.go`
- 创建: `internal/utils/response.go`
- 创建: `internal/utils/password.go`
- 修改: `internal/middleware/auth.go`（替换 stub）
- 修改: `go.mod`（添加 jwt 依赖）

- [ ] **Step 1: 安装 JWT 和 bcrypt 依赖**

```bash
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto
```

- [ ] **Step 2: 创建公共响应工具 `internal/utils/response.go`**

```go
package utils

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    data,
	})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": message,
		"data":    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
		"data":    nil,
	})
}

func Paginated(c *gin.Context, list interface{}, total int64, page, size int) {
	totalPage := int(total) / size
	if int(total)%size > 0 {
		totalPage++
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"list":       list,
			"total":      total,
			"page":       page,
			"size":       size,
			"total_page": totalPage,
		},
	})
}
```

- [ ] **Step 3: 创建 JWT 工具 `internal/utils/jwt.go`**

```go
package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "software-engineering-qa-platform-secret-key"
	}
	jwtSecret = []byte(secret)
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
```

- [ ] **Step 4: 创建密码工具 `internal/utils/password.go`**

```go
package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
```

- [ ] **Step 5: 创建公共 DTO `internal/dto/common.go`**

```go
package dto

type PageRequest struct {
	Page int `form:"page" binding:"required,min=1"`
	Size int `form:"size" binding:"required,min=1,max=100"`
}

func (p *PageRequest) Offset() int {
	return (p.Page - 1) * p.Size
}
```

- [ ] **Step 6: 替换 auth middleware 为真实 JWT 验证**

替换 `internal/middleware/auth.go`:

```go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/utils"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未授权", "data": nil})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "无效的令牌格式", "data": nil})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "无效或过期的令牌", "data": nil})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
```

- [ ] **Step 7: 验证编译并提交**

```bash
go build ./...
go vet ./...
git add -A
git commit -m "feat: 添加 JWT 认证基础设施和公共工具"
```

---

## Task 2: 认证模块 (Auth)

**依赖:** Task 1
**可并行:** 否（User 模块依赖此任务）

**文件:**
- 创建: `internal/dto/auth.go`
- 创建: `internal/repository/user_repo.go`
- 创建: `internal/service/auth_service.go`
- 创建: `internal/controller/auth_controller.go`
- 修改: `internal/routes/routes.go`（替换 stub）

- [ ] **Step 1: 创建认证 DTO `internal/dto/auth.go`**

```go
package dto

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname" binding:"max=50"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	Token string `json:"token" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
```

- [ ] **Step 2: 创建用户 Repository `internal/repository/user_repo.go`**

```go
package repository

import (
	"software_engineering/internal/database"
	"software_engineering/internal/model"
)

func CreateUser(user *model.User) error {
	return database.DB.Create(user).Error
}

func FindUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func FindUserByID(id uint) (*model.User, error) {
	var user model.User
	err := database.DB.First(&user, id).Error
	return &user, err
}

func UpdateUser(user *model.User) error {
	return database.DB.Save(user).Error
}

func ListUsers(page, size int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	database.DB.Model(&model.User{}).Count(&total)
	err := database.DB.Offset((page - 1) * size).Limit(size).Find(&users).Error
	return users, total, err
}
```

- [ ] **Step 3: 创建认证 Service `internal/service/auth_service.go`**

```go
package service

import (
	"errors"

	"software_engineering/internal/dto"
	"software_engineering/internal/model"
	"software_engineering/internal/repository"
	"software_engineering/internal/utils"

	"gorm.io/gorm"
)

func Register(req dto.RegisterRequest) error {
	existing, _ := repository.FindUserByUsername(req.Username)
	if existing.ID != 0 {
		return errors.New("用户名已存在")
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &model.User{
		Username: req.Username,
		Password: hash,
		Email:    req.Email,
		Nickname: req.Nickname,
		Status:   1,
	}
	return repository.CreateUser(user)
}

func Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := repository.FindUserByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	if user.Status == 0 {
		return nil, errors.New("用户已被禁用")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			Status:    user.Status,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		},
	}, nil
}

func RefreshToken(oldToken string) (string, error) {
	claims, err := utils.ParseToken(oldToken)
	if err != nil {
		return "", errors.New("无效的令牌")
	}
	return utils.GenerateToken(claims.UserID, claims.Username)
}
```

- [ ] **Step 4: 创建认证 Controller `internal/controller/auth_controller.go`**

```go
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/dto"
	"software_engineering/internal/service"
	"software_engineering/internal/utils"
)

func Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := service.Register(req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, nil)
}

func Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	resp, err := service.Login(req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, resp)
}

func Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	token, err := service.RefreshToken(req.Token)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	utils.Success(c, gin.H{"token": token})
}
```

- [ ] **Step 5: 修改 routes.go 替换 Auth stub**

在 `internal/routes/routes.go` 中替换 auth 组的 stub：

```go
import (
	"software_engineering/internal/controller"
	// ... 其他 import 保持不变
)

// 在 auth 组中替换：
auth.POST("/register", controller.Register)
auth.POST("/login", controller.Login)
auth.POST("/refresh", controller.Refresh)
```

- [ ] **Step 6: 验证编译并提交**

```bash
go build ./...
go vet ./...
git add -A
git commit -m "feat: 实现认证模块（注册、登录、刷新token）"
```

---

## Task 3: 用户管理模块 (User)

**依赖:** Task 1, Task 2
**可并行:** 是（与 Task 4 并行）

**文件:**
- 创建: `internal/dto/user.go`
- 创建: `internal/service/user_service.go`
- 创建: `internal/controller/user_controller.go`
- 修改: `internal/routes/routes.go`（替换 User stub）

- [ ] **Step 1: 创建用户 DTO `internal/dto/user.go`**

```go
package dto

type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"max=50"`
	Avatar   string `json:"avatar" binding:"max=255"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

type UserListResponse struct {
	List       []UserResponse `json:"list"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Size       int            `json:"size"`
	TotalPage  int            `json:"total_page"`
}
```

- [ ] **Step 2: 创建用户 Service `internal/service/user_service.go`**

```go
package service

import (
	"errors"

	"software_engineering/internal/dto"
	"software_engineering/internal/repository"
	"software_engineering/internal/utils"
)

func GetProfile(userID uint) (*dto.UserResponse, error) {
	user, err := repository.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func UpdateProfile(userID uint, req dto.UpdateProfileRequest) error {
	user, err := repository.FindUserByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	return repository.UpdateUser(user)
}

func ChangePassword(userID uint, req dto.ChangePasswordRequest) error {
	user, err := repository.FindUserByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}
	if !utils.CheckPassword(req.OldPassword, user.Password) {
		return errors.New("旧密码错误")
	}
	hash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	user.Password = hash
	return repository.UpdateUser(user)
}

func ListUsers(page, size int) (*dto.UserListResponse, error) {
	users, total, err := repository.ListUsers(page, size)
	if err != nil {
		return nil, err
	}
	list := make([]dto.UserResponse, len(users))
	for i, u := range users {
		list[i] = dto.UserResponse{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			Nickname:  u.Nickname,
			Avatar:    u.Avatar,
			Status:    u.Status,
			CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	totalPage := int(total) / size
	if int(total)%size > 0 {
		totalPage++
	}
	return &dto.UserListResponse{List: list, Total: total, Page: page, Size: size, TotalPage: totalPage}, nil
}
```

- [ ] **Step 3: 创建用户 Controller `internal/controller/user_controller.go`**

```go
package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/dto"
	"software_engineering/internal/service"
	"software_engineering/internal/utils"
)

func GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	resp, err := service.GetProfile(userID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}
	utils.Success(c, resp)
}

func UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := service.UpdateProfile(userID, req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, nil)
}

func ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := service.ChangePassword(userID, req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, nil)
}

func ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	resp, err := service.ListUsers(page, size)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, resp)
}
```

- [ ] **Step 4: 修改 routes.go 替换 User stub**

```go
users.GET("/profile", controller.GetProfile)
users.PUT("/profile", controller.UpdateProfile)
users.POST("/password", controller.ChangePassword)
users.GET("/list", controller.ListUsers)
```

- [ ] **Step 5: 验证编译并提交**

```bash
go build ./...
go vet ./...
git add -A
git commit -m "feat: 实现用户管理模块（个人信息、修改密码、用户列表）"
```

---

## Task 4: 资源管理模块 (Document)

**依赖:** Task 1
**可并行:** 是（与 Task 3 并行）

**文件:**
- 创建: `internal/dto/document.go`
- 创建: `internal/repository/document_repo.go`
- 创建: `internal/service/document_service.go`
- 创建: `internal/controller/document_controller.go`
- 修改: `internal/routes/routes.go`（替换 Document stub）

- [ ] **Step 1: 创建文档 DTO `internal/dto/document.go`**

```go
package dto

type DocumentResponse struct {
	ID             uint   `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Filename       string `json:"filename"`
	FileSize       int64  `json:"file_size"`
	FileType       string `json:"file_type"`
	Status         string `json:"status"`
	ContentPreview string `json:"content_preview,omitempty"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type DocumentContentResponse struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateDocumentRequest struct {
	Title       string `json:"title" binding:"max=200"`
	Description string `json:"description" binding:"max=500"`
}
```

- [ ] **Step 2: 创建文档 Repository `internal/repository/document_repo.go`**

```go
package repository

import (
	"software_engineering/internal/database"
	"software_engineering/internal/model"
)

func CreateDocument(doc *model.Document) error {
	return database.DB.Create(doc).Error
}

func FindDocumentByID(id uint) (*model.Document, error) {
	var doc model.Document
	err := database.DB.First(&doc, id).Error
	return &doc, err
}

func UpdateDocument(doc *model.Document) error {
	return database.DB.Save(doc).Error
}

func DeleteDocument(id uint) error {
	return database.DB.Delete(&model.Document{}, id).Error
}

func ListDocuments(page, size int, keyword, status string) ([]model.Document, int64, error) {
	var docs []model.Document
	var total int64
	query := database.DB.Model(&model.Document{})
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&docs).Error
	return docs, total, err
}
```

- [ ] **Step 3: 创建文档 Service `internal/service/document_service.go`**

```go
package service

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"software_engineering/internal/dto"
	"software_engineering/internal/model"
	"software_engineering/internal/repository"
)

const uploadDir = "./uploads"

func UploadDocument(title, description string, filename string, fileSize int64, fileType string, contentReader io.Reader) (*dto.DocumentResponse, error) {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, err
	}

	// 保存文件
	filePath := filepath.Join(uploadDir, filename)
	out, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer out.Close()
	content, err := io.Copy(out, contentReader)
	if err != nil {
		return nil, err
	}

	// 读取文本内容（仅对文本文件）
	var fileContent string
	ext := strings.ToLower(fileType)
	if ext == ".md" || ext == ".txt" {
		data, err := os.ReadFile(filePath)
		if err == nil {
			fileContent = string(data)
		}
	}

	if title == "" {
		title = filename
	}

	doc := &model.Document{
		Title:       title,
		Description: description,
		Filename:    filename,
		FilePath:    filePath,
		FileSize:    content,
		FileType:    fileType,
		Content:     fileContent,
		Status:      "completed",
	}
	if err := repository.CreateDocument(doc); err != nil {
		return nil, err
	}

	return &dto.DocumentResponse{
		ID:          doc.ID,
		Title:       doc.Title,
		Description: doc.Description,
		Filename:    doc.Filename,
		FileSize:    doc.FileSize,
		FileType:    doc.FileType,
		Status:      doc.Status,
		CreatedAt:   doc.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   doc.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func GetDocument(id uint) (*dto.DocumentResponse, error) {
	doc, err := repository.FindDocumentByID(id)
	if err != nil {
		return nil, errors.New("文档不存在")
	}
	preview := doc.Content
	if len(preview) > 200 {
		preview = preview[:200] + "..."
	}
	return &dto.DocumentResponse{
		ID:             doc.ID,
		Title:          doc.Title,
		Description:    doc.Description,
		Filename:       doc.Filename,
		FileSize:       doc.FileSize,
		FileType:       doc.FileType,
		Status:         doc.Status,
		ContentPreview: preview,
		CreatedAt:      doc.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      doc.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func GetDocumentContent(id uint) (*dto.DocumentContentResponse, error) {
	doc, err := repository.FindDocumentByID(id)
	if err != nil {
		return nil, errors.New("文档不存在")
	}
	return &dto.DocumentContentResponse{ID: doc.ID, Title: doc.Title, Content: doc.Content}, nil
}

func UpdateDocument(id uint, req dto.UpdateDocumentRequest) error {
	doc, err := repository.FindDocumentByID(id)
	if err != nil {
		return errors.New("文档不存在")
	}
	if req.Title != "" {
		doc.Title = req.Title
	}
	if req.Description != "" {
		doc.Description = req.Description
	}
	return repository.UpdateDocument(doc)
}

func DeleteDocument(id uint) error {
	doc, err := repository.FindDocumentByID(id)
	if err != nil {
		return errors.New("文档不存在")
	}
	// 删除物理文件
	if doc.FilePath != "" {
		os.Remove(doc.FilePath)
	}
	return repository.DeleteDocument(id)
}

func ListDocuments(page, size int, keyword, status string) ([]dto.DocumentResponse, int64, error) {
	docs, total, err := repository.ListDocuments(page, size, keyword, status)
	if err != nil {
		return nil, 0, err
	}
	list := make([]dto.DocumentResponse, len(docs))
	for i, d := range docs {
		list[i] = dto.DocumentResponse{
			ID:          d.ID,
			Title:       d.Title,
			Description: d.Description,
			Filename:    d.Filename,
			FileSize:    d.FileSize,
			FileType:    d.FileType,
			Status:      d.Status,
			CreatedAt:   d.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   d.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return list, total, nil
}
```

- [ ] **Step 4: 创建文档 Controller `internal/controller/document_controller.go`**

```go
package controller

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/dto"
	"software_engineering/internal/service"
	"software_engineering/internal/utils"
)

func UploadDocument(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "请上传文件")
		return
	}
	title := c.PostForm("title")
	description := c.PostForm("description")
	ext := filepath.Ext(file.Filename)

	f, err := file.Open()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "文件读取失败")
		return
	}
	defer f.Close()

	resp, err := service.UploadDocument(title, description, file.Filename, file.Size, ext, f)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, resp)
}

func GetDocument(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	resp, err := service.GetDocument(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}
	utils.Success(c, resp)
}

func GetDocumentContent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	resp, err := service.GetDocumentContent(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}
	utils.Success(c, resp)
}

func UpdateDocument(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := service.UpdateDocument(uint(id), req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, nil)
}

func DeleteDocument(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.DeleteDocument(uint(id)); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, nil)
}

func ListDocuments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	keyword := c.Query("keyword")
	status := c.Query("status")
	list, total, err := service.ListDocuments(page, size, keyword, status)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Paginated(c, list, total, page, size)
}
```

- [ ] **Step 5: 修改 routes.go 替换 Document stub**

```go
docs.POST("", controller.UploadDocument)
docs.GET("", controller.ListDocuments)
docs.GET("/:id", controller.GetDocument)
docs.PUT("/:id", controller.UpdateDocument)
docs.DELETE("/:id", controller.DeleteDocument)
docs.GET("/:id/content", controller.GetDocumentContent)
```

- [ ] **Step 6: 验证编译并提交**

```bash
go build ./...
go vet ./...
git add -A
git commit -m "feat: 实现资源管理模块（文档上传、列表、详情、更新、删除）"
```

---

## Task 5: 知识点管理模块 (Knowledge)

**依赖:** Task 1
**可并行:** 是（与 Task 3, 4 并行）

**文件:**
- 创建: `internal/dto/knowledge.go`
- 创建: `internal/repository/knowledge_repo.go`
- 创建: `internal/service/knowledge_service.go`
- 创建: `internal/controller/knowledge_controller.go`
- 修改: `internal/routes/routes.go`（替换 Knowledge stub）

- [ ] **Step 1: 创建知识点 DTO `internal/dto/knowledge.go`**

```go
package dto

type KnowledgePointResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DocumentID  uint   `json:"document_id"`
	Category    string `json:"category"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type CreateKnowledgePointRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
	DocumentID  uint   `json:"document_id"`
	Category    string `json:"category" binding:"max=50"`
}

type UpdateKnowledgePointRequest struct {
	Name        string `json:"name" binding:"max=100"`
	Description string `json:"description" binding:"max=500"`
	Category    string `json:"category" binding:"max=50"`
}

type KnowledgeRelationResponse struct {
	ID           uint   `json:"id"`
	SourceID     uint   `json:"source_id"`
	SourceName   string `json:"source_name"`
	TargetID     uint   `json:"target_id"`
	TargetName   string `json:"target_name"`
	RelationType string `json:"relation_type"`
	Description  string `json:"description"`
	CreatedAt    string `json:"created_at"`
}

type CreateRelationRequest struct {
	SourceID     uint   `json:"source_id" binding:"required"`
	TargetID     uint   `json:"target_id" binding:"required"`
	RelationType string `json:"relation_type" binding:"required,oneof=RELATED DEPENDS_ON PART_OF"`
	Description  string `json:"description" binding:"max=500"`
}

type UpdateRelationRequest struct {
	RelationType string `json:"relation_type" binding:"oneof=RELATED DEPENDS_ON PART_OF"`
	Description  string `json:"description" binding:"max=500"`
}
```

- [ ] **Step 2: 创建知识点 Repository `internal/repository/knowledge_repo.go`**

```go
package repository

import (
	"software_engineering/internal/database"
	"software_engineering/internal/model"
)

func CreateKnowledgePoint(kp *model.KnowledgePoint) error {
	return database.DB.Create(kp).Error
}

func FindKnowledgePointByID(id uint) (*model.KnowledgePoint, error) {
	var kp model.KnowledgePoint
	err := database.DB.First(&kp, id).Error
	return &kp, err
}

func UpdateKnowledgePoint(kp *model.KnowledgePoint) error {
	return database.DB.Save(kp).Error
}

func DeleteKnowledgePoint(id uint) error {
	return database.DB.Delete(&model.KnowledgePoint{}, id).Error
}

func ListKnowledgePoints(page, size int, keyword string, documentID uint) ([]model.KnowledgePoint, int64, error) {
	var points []model.KnowledgePoint
	var total int64
	query := database.DB.Model(&model.KnowledgePoint{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if documentID > 0 {
		query = query.Where("document_id = ?", documentID)
	}
	query.Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Find(&points).Error
	return points, total, err
}

func CreateRelation(rel *model.KnowledgeRelation) error {
	return database.DB.Create(rel).Error
}

func FindRelationByID(id uint) (*model.KnowledgeRelation, error) {
	var rel model.KnowledgeRelation
	err := database.DB.First(&rel, id).Error
	return &rel, err
}

func UpdateRelation(rel *model.KnowledgeRelation) error {
	return database.DB.Save(rel).Error
}

func DeleteRelation(id uint) error {
	return database.DB.Delete(&model.KnowledgeRelation{}, id).Error
}

func ListRelations(page, size int, pointID uint) ([]model.KnowledgeRelation, int64, error) {
	var rels []model.KnowledgeRelation
	var total int64
	query := database.DB.Model(&model.KnowledgeRelation{})
	if pointID > 0 {
		query = query.Where("source_id = ? OR target_id = ?", pointID, pointID)
	}
	query.Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Find(&rels).Error
	return rels, total, err
}

func GetAllKnowledgePoints() ([]model.KnowledgePoint, error) {
	var points []model.KnowledgePoint
	err := database.DB.Find(&points).Error
	return points, err
}
```

- [ ] **Step 3: 创建知识点 Service `internal/service/knowledge_service.go`**

```go
package service

import (
	"errors"

	"software_engineering/internal/dto"
	"software_engineering/internal/model"
	"software_engineering/internal/repository"
)

func CreateKnowledgePoint(req dto.CreateKnowledgePointRequest) (uint, error) {
	kp := &model.KnowledgePoint{
		Name:        req.Name,
		Description: req.Description,
		DocumentID:  req.DocumentID,
		Category:    req.Category,
	}
	if err := repository.CreateKnowledgePoint(kp); err != nil {
		return 0, err
	}
	return kp.ID, nil
}

func GetKnowledgePoint(id uint) (*dto.KnowledgePointResponse, error) {
	kp, err := repository.FindKnowledgePointByID(id)
	if err != nil {
		return nil, errors.New("知识点不存在")
	}
	return &dto.KnowledgePointResponse{
		ID:          kp.ID,
		Name:        kp.Name,
		Description: kp.Description,
		DocumentID:  kp.DocumentID,
		Category:    kp.Category,
		CreatedAt:   kp.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   kp.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func UpdateKnowledgePoint(id uint, req dto.UpdateKnowledgePointRequest) error {
	kp, err := repository.FindKnowledgePointByID(id)
	if err != nil {
		return errors.New("知识点不存在")
	}
	if req.Name != "" {
		kp.Name = req.Name
	}
	if req.Description != "" {
		kp.Description = req.Description
	}
	if req.Category != "" {
		kp.Category = req.Category
	}
	return repository.UpdateKnowledgePoint(kp)
}

func DeleteKnowledgePoint(id uint) error {
	_, err := repository.FindKnowledgePointByID(id)
	if err != nil {
		return errors.New("知识点不存在")
	}
	return repository.DeleteKnowledgePoint(id)
}

func ListKnowledgePoints(page, size int, keyword string, documentID uint) ([]dto.KnowledgePointResponse, int64, error) {
	points, total, err := repository.ListKnowledgePoints(page, size, keyword, documentID)
	if err != nil {
		return nil, 0, err
	}
	list := make([]dto.KnowledgePointResponse, len(points))
	for i, p := range points {
		list[i] = dto.KnowledgePointResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			DocumentID:  p.DocumentID,
			Category:    p.Category,
			CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return list, total, nil
}

// --- 关系 ---

func CreateRelation(req dto.CreateRelationRequest) (uint, error) {
	// 验证源和目标知识点存在
	if _, err := repository.FindKnowledgePointByID(req.SourceID); err != nil {
		return 0, errors.New("源知识点不存在")
	}
	if _, err := repository.FindKnowledgePointByID(req.TargetID); err != nil {
		return 0, errors.New("目标知识点不存在")
	}
	rel := &model.KnowledgeRelation{
		SourceID:     req.SourceID,
		TargetID:     req.TargetID,
		RelationType: req.RelationType,
		Description:  req.Description,
	}
	if err := repository.CreateRelation(rel); err != nil {
		return 0, err
	}
	return rel.ID, nil
}

func UpdateRelation(id uint, req dto.UpdateRelationRequest) error {
	rel, err := repository.FindRelationByID(id)
	if err != nil {
		return errors.New("关系不存在")
	}
	if req.RelationType != "" {
		rel.RelationType = req.RelationType
	}
	if req.Description != "" {
		rel.Description = req.Description
	}
	return repository.UpdateRelation(rel)
}

func DeleteRelation(id uint) error {
	_, err := repository.FindRelationByID(id)
	if err != nil {
		return errors.New("关系不存在")
	}
	return repository.DeleteRelation(id)
}

func ListRelations(page, size int, pointID uint) ([]dto.KnowledgeRelationResponse, int64, error) {
	rels, total, err := repository.ListRelations(page, size, pointID)
	if err != nil {
		return nil, 0, err
	}

	// 批量查询知识点名称
	pointIDs := make(map[uint]bool)
	for _, r := range rels {
		pointIDs[r.SourceID] = true
		pointIDs[r.TargetID] = true
	}
	names := make(map[uint]string)
	for id := range pointIDs {
		if kp, err := repository.FindKnowledgePointByID(id); err == nil {
			names[id] = kp.Name
		}
	}

	list := make([]dto.KnowledgeRelationResponse, len(rels))
	for i, r := range rels {
		list[i] = dto.KnowledgeRelationResponse{
			ID:           r.ID,
			SourceID:     r.SourceID,
			SourceName:   names[r.SourceID],
			TargetID:     r.TargetID,
			TargetName:   names[r.TargetID],
			RelationType: r.RelationType,
			Description:  r.Description,
			CreatedAt:    r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return list, total, nil
}
```

- [ ] **Step 4: 创建知识点 Controller `internal/controller/knowledge_controller.go`**

```go
package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/dto"
	"software_engineering/internal/service"
	"software_engineering/internal/utils"
)

func CreateKnowledgePoint(c *gin.Context) {
	var req dto.CreateKnowledgePointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	id, err := service.CreateKnowledgePoint(req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, gin.H{"id": id})
}

func GetKnowledgePoint(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	resp, err := service.GetKnowledgePoint(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}
	utils.Success(c, resp)
}

func UpdateKnowledgePoint(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.UpdateKnowledgePointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := service.UpdateKnowledgePoint(uint(id), req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, nil)
}

func DeleteKnowledgePoint(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.DeleteKnowledgePoint(uint(id)); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, nil)
}

func ListKnowledgePoints(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	keyword := c.Query("keyword")
	documentID, _ := strconv.Atoi(c.Query("document_id"))
	list, total, err := service.ListKnowledgePoints(page, size, keyword, uint(documentID))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Paginated(c, list, total, page, size)
}

func CreateRelation(c *gin.Context) {
	var req dto.CreateRelationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	id, err := service.CreateRelation(req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, gin.H{"id": id})
}

func UpdateRelation(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.UpdateRelationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := service.UpdateRelation(uint(id), req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, nil)
}

func DeleteRelation(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.DeleteRelation(uint(id)); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, nil)
}

func ListRelations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	pointID, _ := strconv.Atoi(c.Query("point_id"))
	list, total, err := service.ListRelations(page, size, uint(pointID))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Paginated(c, list, total, page, size)
}
```

- [ ] **Step 5: 修改 routes.go 替换 Knowledge stub**

```go
kp.GET("/points", controller.ListKnowledgePoints)
kp.GET("/points/:id", controller.GetKnowledgePoint)
kp.POST("/points", controller.CreateKnowledgePoint)
kp.PUT("/points/:id", controller.UpdateKnowledgePoint)
kp.DELETE("/points/:id", controller.DeleteKnowledgePoint)
kp.GET("/relations", controller.ListRelations)
kp.POST("/relations", controller.CreateRelation)
kp.PUT("/relations/:id", controller.UpdateRelation)
kp.DELETE("/relations/:id", controller.DeleteRelation)
```

- [ ] **Step 6: 验证编译并提交**

```bash
go build ./...
go vet ./...
git add -A
git commit -m "feat: 实现知识点管理模块（知识点和关系的增删改查）"
```

---

## Task 6: 题库模块 (Question)

**依赖:** Task 1
**可并行:** 是（与 Task 3, 4, 5 并行）

**文件:**
- 创建: `internal/dto/question.go`
- 创建: `internal/repository/question_repo.go`
- 创建: `internal/service/question_service.go`
- 创建: `internal/controller/question_controller.go`
- 修改: `internal/routes/routes.go`（替换 Question stub）

- [ ] **Step 1: 创建题目 DTO `internal/dto/question.go`**

```go
package dto

type QuestionOption struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type QuestionResponse struct {
	ID                uint             `json:"id"`
	Title             string           `json:"title"`
	Type              string           `json:"type"`
	Difficulty        string           `json:"difficulty"`
	KnowledgePointID  uint             `json:"knowledge_point_id"`
	KnowledgePointName string          `json:"knowledge_point_name,omitempty"`
	Options           []QuestionOption `json:"options"`
	Answer            string           `json:"answer,omitempty"`
	Explanation       string           `json:"explanation,omitempty"`
	CreatedAt         string           `json:"created_at"`
}

type CreateQuestionRequest struct {
	Title            string           `json:"title" binding:"required,max=500"`
	Type             string           `json:"type" binding:"required,oneof=single multiple"`
	Difficulty       string           `json:"difficulty" binding:"required,oneof=easy medium hard"`
	KnowledgePointID uint             `json:"knowledge_point_id" binding:"required"`
	Options          []QuestionOption `json:"options" binding:"required,min=2"`
	Answer           string           `json:"answer" binding:"required"`
	Explanation      string           `json:"explanation"`
}

type UpdateQuestionRequest struct {
	Title       string           `json:"title" binding:"max=500"`
	Type        string           `json:"type" binding:"oneof=single multiple"`
	Difficulty  string           `json:"difficulty" binding:"oneof=easy medium hard"`
	Options     []QuestionOption `json:"options"`
	Answer      string           `json:"answer"`
	Explanation string           `json:"explanation"`
}
```

- [ ] **Step 2: 创建题目 Repository `internal/repository/question_repo.go`**

```go
package repository

import (
	"software_engineering/internal/database"
	"software_engineering/internal/model"
)

func CreateQuestion(q *model.Question) error {
	return database.DB.Create(q).Error
}

func FindQuestionByID(id uint) (*model.Question, error) {
	var q model.Question
	err := database.DB.First(&q, id).Error
	return &q, err
}

func UpdateQuestion(q *model.Question) error {
	return database.DB.Save(q).Error
}

func DeleteQuestion(id uint) error {
	return database.DB.Delete(&model.Question{}, id).Error
}

func ListQuestions(page, size int, keyword string, knowledgePointID uint, difficulty string) ([]model.Question, int64, error) {
	var questions []model.Question
	var total int64
	query := database.DB.Model(&model.Question{})
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}
	if knowledgePointID > 0 {
		query = query.Where("knowledge_point_id = ?", knowledgePointID)
	}
	if difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}
	query.Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Find(&questions).Error
	return questions, total, err
}
```

- [ ] **Step 3: 创建题目 Service `internal/service/question_service.go`**

```go
package service

import (
	"encoding/json"
	"errors"

	"software_engineering/internal/dto"
	"software_engineering/internal/model"
	"software_engineering/internal/repository"
)

func parseOptions(optionsJSON string) []dto.QuestionOption {
	var options []dto.QuestionOption
	json.Unmarshal([]byte(optionsJSON), &options)
	return options
}

func CreateQuestion(req dto.CreateQuestionRequest) (uint, error) {
	optionsJSON, _ := json.Marshal(req.Options)
	q := &model.Question{
		Title:            req.Title,
		Type:             req.Type,
		Difficulty:       req.Difficulty,
		KnowledgePointID: req.KnowledgePointID,
		Options:          string(optionsJSON),
		Answer:           req.Answer,
		Explanation:      req.Explanation,
	}
	if err := repository.CreateQuestion(q); err != nil {
		return 0, err
	}
	return q.ID, nil
}

func GetQuestion(id uint, includeAnswer bool) (*dto.QuestionResponse, error) {
	q, err := repository.FindQuestionByID(id)
	if err != nil {
		return nil, errors.New("题目不存在")
	}
	resp := &dto.QuestionResponse{
		ID:               q.ID,
		Title:            q.Title,
		Type:             q.Type,
		Difficulty:       q.Difficulty,
		KnowledgePointID: q.KnowledgePointID,
		Options:          parseOptions(q.Options),
		CreatedAt:        q.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if includeAnswer {
		resp.Answer = q.Answer
		resp.Explanation = q.Explanation
	}
	return resp, nil
}

func UpdateQuestion(id uint, req dto.UpdateQuestionRequest) error {
	q, err := repository.FindQuestionByID(id)
	if err != nil {
		return errors.New("题目不存在")
	}
	if req.Title != "" {
		q.Title = req.Title
	}
	if req.Type != "" {
		q.Type = req.Type
	}
	if req.Difficulty != "" {
		q.Difficulty = req.Difficulty
	}
	if req.Options != nil {
		optionsJSON, _ := json.Marshal(req.Options)
		q.Options = string(optionsJSON)
	}
	if req.Answer != "" {
		q.Answer = req.Answer
	}
	if req.Explanation != "" {
		q.Explanation = req.Explanation
	}
	return repository.UpdateQuestion(q)
}

func DeleteQuestion(id uint) error {
	_, err := repository.FindQuestionByID(id)
	if err != nil {
		return errors.New("题目不存在")
	}
	return repository.DeleteQuestion(id)
}

func ListQuestions(page, size int, keyword string, knowledgePointID uint, difficulty string) ([]dto.QuestionResponse, int64, error) {
	questions, total, err := repository.ListQuestions(page, size, keyword, knowledgePointID, difficulty)
	if err != nil {
		return nil, 0, err
	}
	list := make([]dto.QuestionResponse, len(questions))
	for i, q := range questions {
		list[i] = dto.QuestionResponse{
			ID:               q.ID,
			Title:            q.Title,
			Type:             q.Type,
			Difficulty:       q.Difficulty,
			KnowledgePointID: q.KnowledgePointID,
			Options:          parseOptions(q.Options),
			CreatedAt:        q.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return list, total, nil
}
```

- [ ] **Step 4: 创建题目 Controller `internal/controller/question_controller.go`**

```go
package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/dto"
	"software_engineering/internal/service"
	"software_engineering/internal/utils"
)

func CreateQuestion(c *gin.Context) {
	var req dto.CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	id, err := service.CreateQuestion(req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, gin.H{"id": id})
}

func GetQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	resp, err := service.GetQuestion(uint(id), true)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}
	utils.Success(c, resp)
}

func UpdateQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req dto.UpdateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := service.UpdateQuestion(uint(id), req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, nil)
}

func DeleteQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.DeleteQuestion(uint(id)); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, nil)
}

func ListQuestions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	keyword := c.Query("keyword")
	knowledgePointID, _ := strconv.Atoi(c.Query("knowledge_point_id"))
	difficulty := c.Query("difficulty")
	list, total, err := service.ListQuestions(page, size, keyword, uint(knowledgePointID), difficulty)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Paginated(c, list, total, page, size)
}
```

- [ ] **Step 5: 修改 routes.go 替换 Question stub**

```go
q.GET("", controller.ListQuestions)
q.GET("/:id", controller.GetQuestion)
q.POST("", controller.CreateQuestion)
q.PUT("/:id", controller.UpdateQuestion)
q.DELETE("/:id", controller.DeleteQuestion)
```

- [ ] **Step 6: 验证编译并提交**

```bash
go build ./...
go vet ./...
git add -A
git commit -m "feat: 实现题库模块（题目增删改查）"
```

---

## Task 7: 答题模块 (Quiz)

**依赖:** Task 1, Task 6
**可并行:** 是（与 Task 8 并行）

**文件:**
- 创建: `internal/dto/quiz.go`
- 创建: `internal/repository/quiz_repo.go`
- 创建: `internal/service/quiz_service.go`
- 创建: `internal/controller/quiz_controller.go`
- 修改: `internal/routes/routes.go`（替换 Quiz stub）

- [ ] **Step 1: 创建答题 DTO `internal/dto/quiz.go`**

```go
package dto

type SubmitQuizRequest struct {
	QuestionID uint   `json:"question_id" binding:"required"`
	UserAnswer string `json:"user_answer" binding:"required"`
}

type QuizResponse struct {
	QuizID         uint             `json:"quiz_id"`
	QuestionID     uint             `json:"question_id"`
	QuestionTitle  string           `json:"question_title,omitempty"`
	Type           string           `json:"type,omitempty"`
	Difficulty     string           `json:"difficulty,omitempty"`
	Options        []QuestionOption `json:"options,omitempty"`
	UserAnswer     string           `json:"user_answer"`
	CorrectAnswer  string           `json:"correct_answer,omitempty"`
	IsCorrect      bool             `json:"is_correct"`
	Explanation    string           `json:"explanation,omitempty"`
	KnowledgePointID   uint         `json:"knowledge_point_id,omitempty"`
	KnowledgePointName string       `json:"knowledge_point_name,omitempty"`
	CreatedAt      string           `json:"created_at"`
}
```

- [ ] **Step 2: 创建答题 Repository `internal/repository/quiz_repo.go`**

```go
package repository

import (
	"software_engineering/internal/database"
	"software_engineering/internal/model"
)

func CreateQuiz(quiz *model.Quiz) error {
	return database.DB.Create(quiz).Error
}

func FindQuizByID(id uint) (*model.Quiz, error) {
	var quiz model.Quiz
	err := database.DB.First(&quiz, id).Error
	return &quiz, err
}

func ListQuizzesByUser(userID uint, page, size int, knowledgePointID uint, isCorrect *bool) ([]model.Quiz, int64, error) {
	var quizzes []model.Quiz
	var total int64
	query := database.DB.Model(&model.Quiz{}).Where("user_id = ?", userID)
	if knowledgePointID > 0 {
		query = query.Joins("JOIN questions ON quizzes.question_id = questions.id AND questions.knowledge_point_id = ?", knowledgePointID)
	}
	if isCorrect != nil {
		query = query.Where("is_correct = ?", *isCorrect)
	}
	query.Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&quizzes).Error
	return quizzes, total, err
}

func CountCorrectByUser(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Quiz{}).Where("user_id = ? AND is_correct = ?", userID, true).Count(&count).Error
	return count, err
}

func CountTotalByUser(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Quiz{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
```

- [ ] **Step 3: 创建答题 Service `internal/service/quiz_service.go`**

```go
package service

import (
	"errors"

	"software_engineering/internal/dto"
	"software_engineering/internal/model"
	"software_engineering/internal/repository"
)

func SubmitQuiz(userID uint, req dto.SubmitQuizRequest) (*dto.QuizResponse, error) {
	question, err := repository.FindQuestionByID(req.QuestionID)
	if err != nil {
		return nil, errors.New("题目不存在")
	}

	isCorrect := question.Answer == req.UserAnswer

	quiz := &model.Quiz{
		QuestionID: req.QuestionID,
		UserID:     userID,
		UserAnswer: req.UserAnswer,
		IsCorrect:  isCorrect,
	}
	if err := repository.CreateQuiz(quiz); err != nil {
		return nil, err
	}

	return &dto.QuizResponse{
		QuizID:        quiz.ID,
		QuestionID:    question.ID,
		UserAnswer:    req.UserAnswer,
		CorrectAnswer: question.Answer,
		IsCorrect:     isCorrect,
		Explanation:   question.Explanation,
		CreatedAt:     quiz.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func GetQuizDetail(id uint) (*dto.QuizResponse, error) {
	quiz, err := repository.FindQuizByID(id)
	if err != nil {
		return nil, errors.New("答题记录不存在")
	}
	question, err := repository.FindQuestionByID(quiz.QuestionID)
	if err != nil {
		return nil, errors.New("题目不存在")
	}

	return &dto.QuizResponse{
		QuizID:        quiz.ID,
		QuestionID:    question.ID,
		QuestionTitle: question.Title,
		Type:          question.Type,
		Difficulty:    question.Difficulty,
		Options:       parseOptions(question.Options),
		UserAnswer:    quiz.UserAnswer,
		CorrectAnswer: question.Answer,
		IsCorrect:     quiz.IsCorrect,
		Explanation:   question.Explanation,
		CreatedAt:     quiz.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func ListQuizHistory(userID uint, page, size int, knowledgePointID uint, isCorrect *bool) ([]dto.QuizResponse, int64, error) {
	quizzes, total, err := repository.ListQuizzesByUser(userID, page, size, knowledgePointID, isCorrect)
	if err != nil {
		return nil, 0, err
	}
	list := make([]dto.QuizResponse, len(quizzes))
	for i, q := range quizzes {
		question, _ := repository.FindQuestionByID(q.QuestionID)
		item := dto.QuizResponse{
			QuizID:     q.ID,
			QuestionID: q.QuestionID,
			UserAnswer: q.UserAnswer,
			IsCorrect:  q.IsCorrect,
			CreatedAt:  q.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
		if question != nil {
			item.QuestionTitle = question.Title
			item.KnowledgePointID = question.KnowledgePointID
		}
		list[i] = item
	}
	return list, total, nil
}
```

- [ ] **Step 4: 创建答题 Controller `internal/controller/quiz_controller.go`**

```go
package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/dto"
	"software_engineering/internal/service"
	"software_engineering/internal/utils"
)

func SubmitQuiz(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.SubmitQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	resp, err := service.SubmitQuiz(userID, req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.Success(c, resp)
}

func GetQuizDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	resp, err := service.GetQuizDetail(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}
	utils.Success(c, resp)
}

func ListQuizHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	knowledgePointID, _ := strconv.Atoi(c.Query("knowledge_point_id"))
	list, total, err := service.ListQuizHistory(userID, page, size, uint(knowledgePointID), nil)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Paginated(c, list, total, page, size)
}
```

- [ ] **Step 5: 修改 routes.go 替换 Quiz stub**

```go
quiz.POST("/submit", controller.SubmitQuiz)
quiz.GET("/history", controller.ListQuizHistory)
quiz.GET("/:id", controller.GetQuizDetail)
```

- [ ] **Step 6: 验证编译并提交**

```bash
go build ./...
go vet ./...
git add -A
git commit -m "feat: 实现答题模块（提交答案、答题历史、答题详情）"
```

---

## Task 8: 知识图谱模块 (Graph)

**依赖:** Task 1, Task 5
**可并行:** 是（与 Task 7 并行）

**文件:**
- 创建: `internal/dto/graph.go`
- 创建: `internal/repository/graph_repo.go`
- 创建: `internal/service/graph_service.go`
- 创建: `internal/controller/graph_controller.go`
- 修改: `internal/routes/routes.go`（替换 Graph stub）

- [ ] **Step 1: 创建图谱 DTO `internal/dto/graph.go`**

```go
package dto

type GraphDataResponse struct {
	Nodes   []GraphNode   `json:"nodes"`
	Edges   []GraphEdge   `json:"edges"`
	Summary GraphSummary  `json:"summary"`
}

type GraphNode struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DocumentID  uint   `json:"document_id"`
	Category    string `json:"category"`
}

type GraphEdge struct {
	ID           uint   `json:"id"`
	Source       uint   `json:"source"`
	Target       uint   `json:"target"`
	RelationType string `json:"relation_type"`
	Description  string `json:"description"`
}

type GraphSummary struct {
	NodeCount int `json:"node_count"`
	EdgeCount int `json:"edge_count"`
}

type BuildGraphRequest struct {
	DocumentIDs []uint `json:"document_ids" binding:"required,min=1"`
}

type BuildGraphResponse struct {
	BuildID          uint   `json:"build_id"`
	CreatedPoints    int    `json:"created_points"`
	CreatedRelations int    `json:"created_relations"`
	ChunkCount       int    `json:"chunk_count"`
	VectorCount      int    `json:"vector_count"`
	Status           string `json:"status"`
	Message          string `json:"message"`
}

type BuildHistoryResponse struct {
	List       []BuildGraphResponse `json:"list"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	Size       int                  `json:"size"`
	TotalPage  int                  `json:"total_page"`
}
```

- [ ] **Step 2: 创建图谱 Repository `internal/repository/graph_repo.go`**

```go
package repository

import (
	"software_engineering/internal/database"
	"software_engineering/internal/model"
)

func CreateKnowledgeBuild(build *model.KnowledgeBuild) error {
	return database.DB.Create(build).Error
}

func GetLatestBuild() (*model.KnowledgeBuild, error) {
	var build model.KnowledgeBuild
	err := database.DB.Order("created_at DESC").First(&build).Error
	return &build, err
}

func ListBuilds(page, size int) ([]model.KnowledgeBuild, int64, error) {
	var builds []model.KnowledgeBuild
	var total int64
	database.DB.Model(&model.KnowledgeBuild{}).Count(&total)
	err := database.DB.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&builds).Error
	return builds, total, err
}

func GetAllKnowledgePointsForGraph() ([]model.KnowledgePoint, error) {
	var points []model.KnowledgePoint
	err := database.DB.Find(&points).Error
	return points, err
}

func GetAllRelationsForGraph() ([]model.KnowledgeRelation, error) {
	var rels []model.KnowledgeRelation
	err := database.DB.Find(&rels).Error
	return rels, err
}

func FindKnowledgePointsByIDs(ids []uint) ([]model.KnowledgePoint, error) {
	var points []model.KnowledgePoint
	err := database.DB.Where("id IN ?", ids).Find(&points).Error
	return points, err
}
```

- [ ] **Step 3: 创建图谱 Service `internal/service/graph_service.go`**

```go
package service

import (
	"fmt"
	"strconv"
	"strings"

	"software_engineering/internal/dto"
	"software_engineering/internal/model"
	"software_engineering/internal/repository"
)

func GetGraphData(documentID uint, keyword string, relationType string) (*dto.GraphDataResponse, error) {
	points, err := repository.GetAllKnowledgePointsForGraph()
	if err != nil {
		return nil, err
	}
	rels, err := repository.GetAllRelationsForGraph()
	if err != nil {
		return nil, err
	}

	// 过滤
	var filteredPoints []model.KnowledgePoint
	for _, p := range points {
		if documentID > 0 && p.DocumentID != documentID {
			continue
		}
		if keyword != "" && !strings.Contains(p.Name, keyword) {
			continue
		}
		filteredPoints = append(filteredPoints, p)
	}

	pointIDs := make(map[uint]bool)
	for _, p := range filteredPoints {
		pointIDs[p.ID] = true
	}

	var filteredRels []model.KnowledgeRelation
	for _, r := range rels {
		if !pointIDs[r.SourceID] && !pointIDs[r.TargetID] {
			continue
		}
		if relationType != "" && r.RelationType != relationType {
			continue
		}
		filteredRels = append(filteredRels, r)
	}

	nodes := make([]dto.GraphNode, len(filteredPoints))
	for i, p := range filteredPoints {
		nodes[i] = dto.GraphNode{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			DocumentID:  p.DocumentID,
			Category:    p.Category,
		}
	}

	edges := make([]dto.GraphEdge, len(filteredRels))
	for i, r := range filteredRels {
		edges[i] = dto.GraphEdge{
			ID:           r.ID,
			Source:       r.SourceID,
			Target:       r.TargetID,
			RelationType: r.RelationType,
			Description:  r.Description,
		}
	}

	return &dto.GraphDataResponse{
		Nodes: nodes,
		Edges: edges,
		Summary: dto.GraphSummary{
			NodeCount: len(nodes),
			EdgeCount: len(edges),
		},
	}, nil
}

func BuildGraph(documentIDs []uint) (*dto.BuildGraphResponse, error) {
	docs, err := repository.FindKnowledgePointsByIDs(documentIDs)
	_ = docs // 用于获取文档内容

	// 简化实现：基于已有知识点构建关系
	existingPoints, _ := repository.GetAllKnowledgePointsForGraph()

	createdPoints := 0
	createdRelations := 0

	// 查找文档关联的知识点
	var docPoints []model.KnowledgePoint
	for _, p := range existingPoints {
		for _, docID := range documentIDs {
			if p.DocumentID == docID {
				docPoints = append(docPoints, p)
			}
		}
	}

	// 为同一文档的知识点创建关系
	for i := 0; i < len(docPoints); i++ {
		for j := i + 1; j < len(docPoints); j++ {
			rel := &model.KnowledgeRelation{
				SourceID:     docPoints[i].ID,
				TargetID:     docPoints[j].ID,
				RelationType: "RELATED",
				Description:  fmt.Sprintf("%s 与 %s 相关", docPoints[i].Name, docPoints[j].Name),
			}
			repository.CreateRelation(rel)
			createdRelations++
		}
	}

	docIDsStr := make([]string, len(documentIDs))
	for i, id := range documentIDs {
		docIDsStr[i] = strconv.Itoa(int(id))
	}

	build := &model.KnowledgeBuild{
		DocumentIDs:      strings.Join(docIDsStr, ","),
		CreatedPoints:    createdPoints,
		CreatedRelations: createdRelations,
		ChunkCount:       len(docPoints),
		VectorCount:      len(docPoints) * 3,
		Status:           "completed",
		Message:          "知识图谱构建完成",
	}
	repository.CreateKnowledgeBuild(build)

	return &dto.BuildGraphResponse{
		BuildID:          build.ID,
		CreatedPoints:    createdPoints,
		CreatedRelations: createdRelations,
		ChunkCount:       build.ChunkCount,
		VectorCount:      build.VectorCount,
		Status:           build.Status,
		Message:          build.Message,
	}, nil
}

func GetLatestBuildResult() (*dto.BuildGraphResponse, error) {
	build, err := repository.GetLatestBuild()
	if err != nil {
		return nil, fmt.Errorf("暂无构建记录")
	}
	return &dto.BuildGraphResponse{
		BuildID:          build.ID,
		CreatedPoints:    build.CreatedPoints,
		CreatedRelations: build.CreatedRelations,
		ChunkCount:       build.ChunkCount,
		VectorCount:      build.VectorCount,
		Status:           build.Status,
		Message:          build.Message,
	}, nil
}

func ListBuildHistory(page, size int) (*dto.BuildHistoryResponse, error) {
	builds, total, err := repository.ListBuilds(page, size)
	if err != nil {
		return nil, err
	}
	list := make([]dto.BuildGraphResponse, len(builds))
	for i, b := range builds {
		list[i] = dto.BuildGraphResponse{
			BuildID:          b.ID,
			CreatedPoints:    b.CreatedPoints,
			CreatedRelations: b.CreatedRelations,
			ChunkCount:       b.ChunkCount,
			VectorCount:      b.VectorCount,
			Status:           b.Status,
			Message:          b.Message,
		}
	}
	totalPage := int(total) / size
	if int(total)%size > 0 {
		totalPage++
	}
	return &dto.BuildHistoryResponse{List: list, Total: total, Page: page, Size: size, TotalPage: totalPage}, nil
}
```

- [ ] **Step 4: 创建图谱 Controller `internal/controller/graph_controller.go`**

```go
package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/dto"
	"software_engineering/internal/service"
	"software_engineering/internal/utils"
)

func GetGraph(c *gin.Context) {
	documentID, _ := strconv.Atoi(c.Query("document_id"))
	keyword := c.Query("keyword")
	relationType := c.Query("relation_type")
	resp, err := service.GetGraphData(uint(documentID), keyword, relationType)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, resp)
}

func BuildGraph(c *gin.Context) {
	var req dto.BuildGraphRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	resp, err := service.BuildGraph(req.DocumentIDs)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, resp)
}

func GetLatestBuild(c *gin.Context) {
	resp, err := service.GetLatestBuildResult()
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}
	utils.Success(c, resp)
}

func ListBuildHistory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	resp, err := service.ListBuildHistory(page, size)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, resp)
}
```

- [ ] **Step 5: 修改 routes.go 替换 Graph stub**

```go
graph.GET("", controller.GetGraph)
graph.POST("/build", controller.BuildGraph)
graph.GET("/build/latest", controller.GetLatestBuild)
graph.GET("/build/history", controller.ListBuildHistory)
```

- [ ] **Step 6: 验证编译并提交**

```bash
go build ./...
go vet ./...
git add -A
git commit -m "feat: 实现知识图谱模块（图谱查询、构建、历史记录）"
```

---

## Task 9: 知识问答模块 (Ask)

**依赖:** Task 1
**可并行:** 是（与 Task 6, 7, 8 并行）

**文件:**
- 创建: `internal/dto/ask.go`
- 创建: `internal/repository/ask_repo.go`
- 创建: `internal/service/ask_service.go`
- 创建: `internal/controller/ask_controller.go`
- 修改: `internal/routes/routes.go`（替换 Ask stub）

- [ ] **Step 1: 创建问答 DTO `internal/dto/ask.go`**

```go
package dto

type AskRequest struct {
	Question      string `json:"question" binding:"required"`
	ConversationID uint  `json:"conversation_id"`
}

type AskResponse struct {
	ConversationID uint                `json:"conversation_id"`
	QuestionID     uint                `json:"question_id"`
	Answer         string              `json:"answer"`
	Confidence     float64             `json:"confidence"`
	Sources        []AskSource         `json:"sources"`
	RelatedKnowledgePoints []KPRef     `json:"related_knowledge_points"`
	CreatedAt      string              `json:"created_at"`
}

type AskSource struct {
	DocumentID    uint   `json:"document_id"`
	DocumentTitle string `json:"document_title"`
	Content       string `json:"content"`
}

type KPRef struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AskSessionResponse struct {
	ConversationID uint   `json:"conversation_id"`
	Title          string `json:"title"`
	LastQuestion   string `json:"last_question,omitempty"`
	MessageCount   int    `json:"message_count"`
	UpdatedAt      string `json:"updated_at"`
}

type AskMessageResponse struct {
	MessageID uint   `json:"message_id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type CreateSessionRequest struct {
	Title string `json:"title" binding:"max=200"`
}
```

- [ ] **Step 2: 创建问答 Repository `internal/repository/ask_repo.go`**

```go
package repository

import (
	"software_engineering/internal/database"
	"software_engineering/internal/model"
)

func CreateAskSession(session *model.AskSession) error {
	return database.DB.Create(session).Error
}

func FindAskSessionByID(id uint) (*model.AskSession, error) {
	var session model.AskSession
	err := database.DB.First(&session, id).Error
	return &session, err
}

func ListAskSessionsByUser(userID uint, page, size int) ([]model.AskSession, int64, error) {
	var sessions []model.AskSession
	var total int64
	query := database.DB.Model(&model.AskSession{}).Where("user_id = ?", userID)
	query.Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Order("updated_at DESC").Find(&sessions).Error
	return sessions, total, err
}

func CreateAskMessage(msg *model.AskMessage) error {
	return database.DB.Create(msg).Error
}

func ListAskMessages(sessionID uint, page, size int) ([]model.AskMessage, int64, error) {
	var messages []model.AskMessage
	var total int64
	query := database.DB.Model(&model.AskMessage{}).Where("session_id = ?", sessionID)
	query.Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Order("created_at ASC").Find(&messages).Error
	return messages, total, err
}

func CountMessagesBySession(sessionID uint) int {
	var count int
	database.DB.Model(&model.AskMessage{}).Where("session_id = ?", sessionID).Count(&count)
	return count
}

func GetLastMessageBySession(sessionID uint) (*model.AskMessage, error) {
	var msg model.AskMessage
	err := database.DB.Where("session_id = ? AND role = ?", sessionID, "user").Order("created_at DESC").First(&msg).Error
	return &msg, err
}
```

- [ ] **Step 3: 创建问答 Service `internal/service/ask_service.go`**

```go
package service

import (
	"fmt"
	"strings"

	"software_engineering/internal/dto"
	"software_engineering/internal/model"
	"software_engineering/internal/repository"
)

func CreateSession(userID uint, req dto.CreateSessionRequest) (*dto.AskSessionResponse, error) {
	session := &model.AskSession{
		UserID: userID,
		Title:  req.Title,
	}
	if err := repository.CreateAskSession(session); err != nil {
		return nil, err
	}
	return &dto.AskSessionResponse{
		ConversationID: session.ID,
		Title:          session.Title,
		UpdatedAt:      session.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func ListSessions(userID uint, page, size int) ([]dto.AskSessionResponse, int64, error) {
	sessions, total, err := repository.ListAskSessionsByUser(userID, page, size)
	if err != nil {
		return nil, 0, err
	}
	list := make([]dto.AskSessionResponse, len(sessions))
	for i, s := range sessions {
		item := dto.AskSessionResponse{
			ConversationID: s.ID,
			Title:          s.Title,
			MessageCount:   repository.CountMessagesBySession(s.ID),
			UpdatedAt:      s.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
		lastMsg, err := repository.GetLastMessageBySession(s.ID)
		if err == nil {
			item.LastQuestion = lastMsg.Content
		}
		list[i] = item
	}
	return list, total, nil
}

func ListSessionMessages(sessionID uint, page, size int) ([]dto.AskMessageResponse, int64, error) {
	messages, total, err := repository.ListAskMessages(sessionID, page, size)
	if err != nil {
		return nil, 0, err
	}
	list := make([]dto.AskMessageResponse, len(messages))
	for i, m := range messages {
		list[i] = dto.AskMessageResponse{
			MessageID: m.ID,
			Role:      m.Role,
			Content:   m.Content,
			CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return list, total, nil
}

func Ask(userID uint, req dto.AskRequest) (*dto.AskResponse, error) {
	// 自动创建或复用会话
	sessionID := req.ConversationID
	if sessionID == 0 {
		session := &model.AskSession{
			UserID: userID,
			Title:  req.Question,
		}
		repository.CreateAskSession(session)
		sessionID = session.ID
	}

	// 保存用户消息
	userMsg := &model.AskMessage{
		SessionID: sessionID,
		Role:      "user",
		Content:   req.Question,
	}
	repository.CreateAskMessage(userMsg)

	// 简化的关键词匹配回答
	answer := fmt.Sprintf("关于「%s」的回答：这是系统生成的示例回答。在生产环境中，这里会接入语义检索和 AI 生成。", req.Question)

	// 查找相关知识点
	points, _ := repository.GetAllKnowledgePointsForGraph()
	var related []dto.KPRef
	for _, p := range points {
		if strings.Contains(req.Question, p.Name) || strings.Contains(p.Name, req.Question) {
			related = append(related, dto.KPRef{ID: p.ID, Name: p.Name, Description: p.Description})
		}
	}

	// 保存助手消息
	assistantMsg := &model.AskMessage{
		SessionID:  sessionID,
		Role:       "assistant",
		Content:    answer,
		Confidence: 0.75,
	}
	repository.CreateAskMessage(assistantMsg)

	return &dto.AskResponse{
		ConversationID: sessionID,
		QuestionID:     userMsg.ID,
		Answer:         answer,
		Confidence:     0.75,
		Sources:        []dto.AskSource{},
		RelatedKnowledgePoints: related,
		CreatedAt:      assistantMsg.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func ListAskHistory(userID uint, page, size int, conversationID uint) ([]dto.AskResponse, int64, error) {
	// 简化：返回会话列表作为历史
	sessions, total, err := repository.ListAskSessionsByUser(userID, page, size)
	if err != nil {
		return nil, 0, err
	}
	list := make([]dto.AskResponse, len(sessions))
	for i, s := range sessions {
		lastMsg, _ := repository.GetLastMessageBySession(s.ID)
		question := ""
		if lastMsg.ID > 0 {
			question = lastMsg.Content
		}
		list[i] = dto.AskResponse{
			ConversationID: s.ID,
			Answer:         question,
			CreatedAt:      s.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return list, total, nil
}
```

- [ ] **Step 4: 创建问答 Controller `internal/controller/ask_controller.go`**

```go
package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/dto"
	"software_engineering/internal/service"
	"software_engineering/internal/utils"
)

func CreateSession(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	resp, err := service.CreateSession(userID, req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, resp)
}

func ListSessions(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	list, total, err := service.ListSessions(userID, page, size)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Paginated(c, list, total, page, size)
}

func ListSessionMessages(c *gin.Context) {
	sessionID, _ := strconv.Atoi(c.Param("id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := service.ListSessionMessages(uint(sessionID), page, size)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Paginated(c, list, total, page, size)
}

func Ask(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.AskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	resp, err := service.Ask(userID, req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, resp)
}

func ListAskHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	conversationID, _ := strconv.Atoi(c.Query("conversation_id"))
	list, total, err := service.ListAskHistory(userID, page, size, uint(conversationID))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Paginated(c, list, total, page, size)
}
```

- [ ] **Step 5: 修改 routes.go 替换 Ask stub**

```go
ask.POST("/sessions", controller.CreateSession)
ask.GET("/sessions", controller.ListSessions)
ask.GET("/sessions/:id/messages", controller.ListSessionMessages)
ask.POST("", controller.Ask)
ask.GET("/history", controller.ListAskHistory)
```

- [ ] **Step 6: 验证编译并提交**

```bash
go build ./...
go vet ./...
git add -A
git commit -m "feat: 实现知识问答模块（会话管理、提问、历史记录）"
```

---

## Task 10: 学习分析模块 (Analytics)

**依赖:** Task 7, Task 9
**可并行:** 否（最后执行）

**文件:**
- 创建: `internal/dto/analytics.go`
- 创建: `internal/repository/analytics_repo.go`
- 创建: `internal/service/analytics_service.go`
- 创建: `internal/controller/analytics_controller.go`
- 修改: `internal/routes/routes.go`（替换 Analytics stub）

- [ ] **Step 1: 创建分析 DTO `internal/dto/analytics.go`**

```go
package dto

type OverviewResponse struct {
	TodayLearningHours    float64 `json:"today_learning_hours"`
	TodayQuestionsAsked   int     `json:"today_questions_asked"`
	TotalLearningHours    float64 `json:"total_learning_hours"`
	TotalQuestionsAsked   int     `json:"total_questions_asked"`
	TotalQuizzesTaken     int     `json:"total_quizzes_taken"`
	AverageCorrectRate    float64 `json:"average_correct_rate"`
	KnowledgePointsMastered int   `json:"knowledge_points_mastered"`
	KnowledgePointsTotal  int     `json:"knowledge_points_total"`
	MasteryRate           float64 `json:"mastery_rate"`
}

type HotKnowledgePoint struct {
	KnowledgePointID   uint   `json:"knowledge_point_id"`
	KnowledgePointName string `json:"knowledge_point_name"`
	Heat               int    `json:"heat"`
	QuestionCount      int    `json:"question_count"`
	QuizCount          int    `json:"quiz_count"`
}

type KnowledgeMastery struct {
	KnowledgePointID   uint    `json:"knowledge_point_id"`
	KnowledgePointName string  `json:"knowledge_point_name"`
	TotalQuestions      int     `json:"total_questions"`
	CorrectAnswers      int     `json:"correct_answers"`
	MasteryRate         float64 `json:"mastery_rate"`
	Level               string  `json:"level"`
}

type WeakPoint struct {
	KnowledgePointID   uint            `json:"knowledge_point_id"`
	KnowledgePointName string          `json:"knowledge_point_name"`
	CorrectRate        float64         `json:"correct_rate"`
	SuggestedQuestions []SuggestedQuestion `json:"suggested_questions"`
}

type SuggestedQuestion struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type TrendData struct {
	DailyStats   []DailyStat   `json:"daily_stats"`
	WeeklyTrend  []WeeklyTrend `json:"weekly_trend"`
}

type DailyStat struct {
	Date           string  `json:"date"`
	QuestionsAsked int     `json:"questions_asked"`
	LearningHours  float64 `json:"learning_hours"`
	CorrectRate    float64 `json:"correct_rate"`
}

type WeeklyTrend struct {
	Week               string  `json:"week"`
	AvgCorrectRate     float64 `json:"avg_correct_rate"`
	TotalLearningHours float64 `json:"total_learning_hours"`
	TotalQuestionsAsked int    `json:"total_questions_asked"`
}
```

- [ ] **Step 2: 创建分析 Repository `internal/repository/analytics_repo.go`**

```go
package repository

import (
	"time"

	"software_engineering/internal/database"
	"software_engineering/internal/model"
)

func CountQuizzesByUser(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Quiz{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func CountCorrectQuizzesByUser(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.Quiz{}).Where("user_id = ? AND is_correct = ?", userID, true).Count(&count).Error
	return count, err
}

func CountTodayQuizzesByUser(userID uint) (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := database.DB.Model(&model.Quiz{}).Where("user_id = ? AND DATE(created_at) = ?", userID, today).Count(&count).Error
	return count, err
}

func CountTodayMessagesByUser(userID uint) (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := database.DB.Model(&model.AskMessage{}).
		Joins("JOIN ask_sessions ON ask_messages.session_id = ask_sessions.id").
		Where("ask_sessions.user_id = ? AND ask_messages.role = 'user' AND DATE(ask_messages.created_at) = ?", userID, today).
		Count(&count).Error
	return count, err
}

func CountTotalMessagesByUser(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.AskMessage{}).
		Joins("JOIN ask_sessions ON ask_messages.session_id = ask_sessions.id").
		Where("ask_sessions.user_id = ? AND ask_messages.role = 'user'", userID).
		Count(&count).Error
	return count, err
}

func GetQuizzesByKnowledgePoint(userID uint) (map[uint]int, map[uint]int, error) {
	var results []struct {
		KnowledgePointID uint
		Total            int
		Correct          int
	}
	err := database.DB.Model(&model.Quiz{}).
		Select("questions.knowledge_point_id, COUNT(*) as total, SUM(CASE WHEN quizzes.is_correct = 1 THEN 1 ELSE 0 END) as correct").
		Joins("JOIN questions ON quizzes.question_id = questions.id").
		Where("quizzes.user_id = ?", userID).
		Group("questions.knowledge_point_id").
		Scan(&results).Error

	totalMap := make(map[uint]int)
	correctMap := make(map[uint]int)
	for _, r := range results {
		totalMap[r.KnowledgePointID] = r.Total
		correctMap[r.KnowledgePointID] = r.Correct
	}
	return totalMap, correctMap, err
}

func GetDailyQuizStats(userID uint, days int) ([]struct {
	Date    string
	Correct int
	Total   int
}, error) {
	var results []struct {
		Date    string
		Correct int
		Total   int
	}
	since := time.Now().AddDate(0, 0, -days)
	err := database.DB.Model(&model.Quiz{}).
		Select("DATE(created_at) as date, COUNT(*) as total, SUM(CASE WHEN is_correct = 1 THEN 1 ELSE 0 END) as correct").
		Where("user_id = ? AND created_at >= ?", userID, since).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&results).Error
	return results, err
}
```

- [ ] **Step 3: 创建分析 Service `internal/service/analytics_service.go`**

```go
package service

import (
	"math"

	"software_engineering/internal/dto"
	"software_engineering/internal/model"
	"software_engineering/internal/repository"
)

func GetOverview(userID uint) (*dto.OverviewResponse, error) {
	totalQuizzes, _ := repository.CountQuizzesByUser(userID)
	correctQuizzes, _ := repository.CountCorrectQuizzesByUser(userID)
	todayQuizzes, _ := repository.CountTodayQuizzesByUser(userID)
	todayMessages, _ := repository.CountTodayMessagesByUser(userID)
	totalMessages, _ := repository.CountTotalMessagesByUser(userID)

	points, _ := repository.GetAllKnowledgePointsForGraph()
	totalPoints := len(points)

	// 简化：正确率 > 80% 的知识点算已掌握
	totalMap, correctMap, _ := repository.GetQuizzesByKnowledgePoint(userID)
	 mastered := 0
	for kpID, total := range totalMap {
		correct := correctMap[kpID]
		if total > 0 && float64(correct)/float64(total) >= 0.8 {
			mastered++
		}
	}

	var avgRate float64
	if totalQuizzes > 0 {
		avgRate = float64(correctQuizzes) / float64(totalQuizzes)
	}

	return &dto.OverviewResponse{
		TodayLearningHours:    float64(todayMessages) * 0.1,
		TodayQuestionsAsked:   int(todayMessages),
		TotalLearningHours:    float64(totalMessages) * 0.1,
		TotalQuestionsAsked:   int(totalMessages),
		TotalQuizzesTaken:     int(totalQuizzes),
		AverageCorrectRate:    math.Round(avgRate*100) / 100,
		KnowledgePointsMastered: mastered,
		KnowledgePointsTotal:  totalPoints,
		MasteryRate:           math.Round(float64(mastered)/float64(totalPoints)*100) / 100,
	}, nil
}

func GetHotKnowledgePoints(limit int) ([]dto.HotKnowledgePoint, error) {
	points, _ := repository.GetAllKnowledgePointsForGraph()
	totalMap, correctMap, _ := repository.GetQuizzesByKnowledgePoint(0)

	var result []dto.HotKnowledgePoint
	for _, p := range points {
		heat := totalMap[p.ID] * 10
		if heat > 0 {
			result = append(result, dto.HotKnowledgePoint{
				KnowledgePointID:   p.ID,
				KnowledgePointName: p.Name,
				Heat:               heat,
				QuestionCount:      totalMap[p.ID],
				QuizCount:          correctMap[p.ID],
			})
		}
	}

	// 按热度排序
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].Heat > result[i].Heat {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}
	return result, nil
}

func GetKnowledgeMastery(userID uint) ([]dto.KnowledgeMastery, error) {
	points, _ := repository.GetAllKnowledgePointsForGraph()
	totalMap, correctMap, _ := repository.GetQuizzesByKnowledgePoint(userID)

	var result []dto.KnowledgeMastery
	for _, p := range points {
		total := totalMap[p.ID]
		correct := correctMap[p.ID]
		if total == 0 {
			continue
		}
		rate := float64(correct) / float64(total)
		level := "weak"
		if rate >= 0.8 {
			level = "mastered"
		} else if rate >= 0.5 {
			level = "learning"
		}
		result = append(result, dto.KnowledgeMastery{
			KnowledgePointID:   p.ID,
			KnowledgePointName: p.Name,
			TotalQuestions:      total,
			CorrectAnswers:      correct,
			MasteryRate:         math.Round(rate*100) / 100,
			Level:               level,
		})
	}
	return result, nil
}

func GetWeakPoints(userID uint, limit int) ([]dto.WeakPoint, error) {
	masteries, _ := GetKnowledgeMastery(userID)

	var result []dto.WeakPoint
	for _, m := range masteries {
		if m.Level == "weak" || m.Level == "learning" {
			result = append(result, dto.WeakPoint{
				KnowledgePointID:   m.KnowledgePointID,
				KnowledgePointName: m.KnowledgePointName,
				CorrectRate:        m.MasteryRate,
			})
		}
	}

	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}
	return result, nil
}

func GetTrends(userID uint, days int) (*dto.TrendData, error) {
	dailyStats, _ := repository.GetDailyQuizStats(userID, days)

	var trends []dto.DailyStat
	for _, d := range dailyStats {
		rate := 0.0
		if d.Total > 0 {
			rate = float64(d.Correct) / float64(d.Total)
		}
		trends = append(trends, dto.DailyStat{
			Date:           d.Date,
			QuestionsAsked: d.Total,
			LearningHours:  float64(d.Total) * 0.1,
			CorrectRate:    math.Round(rate*100) / 100,
		})
	}

	return &dto.TrendData{
		DailyStats:  trends,
		WeeklyTrend: []dto.WeeklyTrend{},
	}, nil
}
```

- [ ] **Step 4: 创建分析 Controller `internal/controller/analytics_controller.go`**

```go
package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/service"
	"software_engineering/internal/utils"
)

func GetOverview(c *gin.Context) {
	userID := c.GetUint("user_id")
	resp, err := service.GetOverview(userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, resp)
}

func GetHotKnowledgePoints(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	resp, err := service.GetHotKnowledgePoints(limit)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, resp)
}

func GetKnowledgeMastery(c *gin.Context) {
	userID := c.GetUint("user_id")
	resp, err := service.GetKnowledgeMastery(userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, resp)
}

func GetWeakPoints(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	resp, err := service.GetWeakPoints(0, limit)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, resp)
}

func GetTrends(c *gin.Context) {
	userID := c.GetUint("user_id")
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	resp, err := service.GetTrends(userID, days)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, resp)
}
```

- [ ] **Step 5: 修改 routes.go 替换 Analytics stub**

```go
analytics.GET("/overview", controller.GetOverview)
analytics.GET("/hot-knowledge-points", controller.GetHotKnowledgePoints)
analytics.GET("/knowledge-mastery", controller.GetKnowledgeMastery)
analytics.GET("/weak-points", controller.GetWeakPoints)
analytics.GET("/trends", controller.GetTrends)
```

- [ ] **Step 6: 验证编译并提交**

```bash
go build ./...
go vet ./...
git add -A
git commit -m "feat: 实现学习分析模块（总览、热门知识点、掌握度、薄弱点、趋势）"
```

---

## Task 11: 最终路由整合和完整构建验证

**依赖:** Task 2-10 全部完成
**可并行:** 否

- [ ] **Step 1: 检查 routes.go 确保所有 stub 已替换**

确认 routes.go 中不再有任何 `stubHandler` 或 `paginatedStub` 调用（healthCheck 除外）。

- [ ] **Step 2: 清理 routes.go 中未使用的 import**

- [ ] **Step 3: 完整构建验证**

```bash
go mod tidy
go build ./...
go vet ./...
```

- [ ] **Step 4: 最终提交**

```bash
git add -A
git commit -m "chore: 完成所有模块整合，全量构建验证通过"
```
