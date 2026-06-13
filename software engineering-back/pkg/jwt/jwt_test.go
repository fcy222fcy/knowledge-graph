package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name     string
		userID   uint
		username string
		wantErr  bool
	}{
		{
			name:     "正常生成Token",
			userID:   1,
			username: "testuser",
			wantErr:  false,
		},
		{
			name:     "用户ID为0",
			userID:   0,
			username: "testuser",
			wantErr:  false,
		},
		{
			name:     "空用户名",
			userID:   1,
			username: "",
			wantErr:  false,
		},
		{
			name:     "长用户名",
			userID:   100,
			username: "very-long-username-for-testing-purpose",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && token == "" {
				t.Error("GenerateToken() 返回空字符串")
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	// 先生成一个有效 token
	validToken, err := GenerateToken(1, "testuser")
	if err != nil {
		t.Fatalf("GenerateToken() 失败: %v", err)
	}

	tests := []struct {
		name        string
		tokenString string
		wantErr     bool
		wantUserID  uint
		wantUser    string
	}{
		{
			name:        "有效Token",
			tokenString: validToken,
			wantErr:     false,
			wantUserID:  1,
			wantUser:    "testuser",
		},
		{
			name:        "空字符串",
			tokenString: "",
			wantErr:     true,
		},
		{
			name:        "无效格式",
			tokenString: "invalid.token.here",
			wantErr:     true,
		},
		{
			name:        "篡改签名",
			tokenString: validToken + "tampered",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ParseToken(tt.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if claims.UserID != tt.wantUserID {
					t.Errorf("ParseToken() UserID = %v, want %v", claims.UserID, tt.wantUserID)
				}
				if claims.Username != tt.wantUser {
					t.Errorf("ParseToken() Username = %v, want %v", claims.Username, tt.wantUser)
				}
			}
		})
	}
}

func TestTokenExpiration(t *testing.T) {
	// 创建一个已过期的 token
	claims := Claims{
		UserID:   1,
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // 1小时前过期
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		t.Fatalf("签名失败: %v", err)
	}

	_, err = ParseToken(tokenString)
	if err == nil {
		t.Error("ParseToken() 应该返回过期错误，但没有")
	}
}

func TestTokenClaims(t *testing.T) {
	userID := uint(42)
	username := "claims-test-user"

	token, err := GenerateToken(userID, username)
	if err != nil {
		t.Fatalf("GenerateToken() 失败: %v", err)
	}

	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken() 失败: %v", err)
	}

	// 验证 Claims 字段
	if claims.UserID != userID {
		t.Errorf("UserID = %v, want %v", claims.UserID, userID)
	}
	if claims.Username != username {
		t.Errorf("Username = %v, want %v", claims.Username, username)
	}

	// 验证时间字段
	if claims.ExpiresAt == nil {
		t.Error("ExpiresAt 不应为 nil")
	}
	if claims.IssuedAt == nil {
		t.Error("IssuedAt 不应为 nil")
	}
	if claims.ExpiresAt.Time.Before(time.Now()) {
		t.Error("ExpiresAt 应该在未来")
	}
}
