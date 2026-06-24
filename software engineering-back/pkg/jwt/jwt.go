package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret JWT 签名密钥，从环境变量读取，未配置时使用默认值
var jwtSecret []byte

// init 初始化 JWT 密钥
func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "software-engineering-qa-platform-secret-key"
	}
	jwtSecret = []byte(secret)
}

// Claims JWT 声明结构体，包含用户 ID、用户名和角色
type Claims struct {
	UserID   uint   `json:"user_id"`   // 用户ID
	Username string `json:"username"`  // 用户名
	Role     string `json:"role"`      // 角色: admin, teacher, student
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token，有效期 24 小时
func GenerateToken(userID uint, username, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析并验证 JWT Token，返回声明信息
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
