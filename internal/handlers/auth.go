package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"software_engineering/internal/database"
	"software_engineering/internal/model"
)

var jwtSecret = []byte("software-engineering-secret-key-2024")

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname" binding:"max=50"`
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// 查找用户
	var user model.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或密码错误",
			"data":    nil,
		})
		return
	}

	// 检查用户状态
	if user.Status == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "用户已被禁用",
			"data":    nil,
		})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或密码错误",
			"data":    nil,
		})
		return
	}

	// 生成 JWT
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成令牌失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data": gin.H{
			"token": tokenString,
			"user": gin.H{
				"id":         user.ID,
				"username":   user.Username,
				"email":      user.Email,
				"nickname":   user.Nickname,
				"avatar":     user.Avatar,
				"status":     user.Status,
				"created_at": user.CreatedAt.Format("2006-01-02T15:04:05Z"),
				"updated_at": user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			},
		},
	})
}

// Register 用户注册
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// 检查用户名是否已存在
	var existingUser model.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"code":    409,
			"message": "用户名已存在",
			"data":    nil,
		})
		return
	}

	// 检查邮箱是否已存在
	var existingEmail model.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingEmail).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"code":    409,
			"message": "邮箱已存在",
			"data":    nil,
		})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "密码加密失败",
			"data":    nil,
		})
		return
	}

	// 创建用户
	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Nickname: req.Nickname,
		Status:   1,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建用户失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "注册成功",
		"data":    nil,
	})
}

// RefreshToken 刷新令牌
func RefreshToken(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// 解析旧令牌
	token, err := jwt.ParseWithClaims(req.Token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "无效的令牌",
			"data":    nil,
		})
		return
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "无效的令牌",
			"data":    nil,
		})
		return
	}

	// 验证用户是否存在
	var user model.User
	if err := database.DB.First(&user, claims.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户不存在",
			"data":    nil,
		})
		return
	}

	if user.Status == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "用户已被禁用",
			"data":    nil,
		})
		return
	}

	// 生成新令牌
	newClaims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	newTokenString, err := newToken.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成令牌失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "刷新成功",
		"data": gin.H{
			"token": newTokenString,
		},
	})
}
