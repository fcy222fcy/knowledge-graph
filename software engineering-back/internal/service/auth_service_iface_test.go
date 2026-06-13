package service

import (
	"testing"

	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/model/entity"
	"software_engineering/internal/repository"
	"software_engineering/pkg/bcrypt" //nolint:importas
	"software_engineering/pkg/jwt"
)

func setupTestService() (*AuthServiceImpl, *repository.MockUserRepository) {
	mockRepo := repository.NewMockUserRepository()
	svc := NewAuthService(mockRepo)
	return svc, mockRepo
}

func TestAuthServiceImpl_Register_Success(t *testing.T) {
	svc, _ := setupTestService()

	req := request.RegisterRequest{
		Username: "newuser",
		Password: "password123",
		Email:    "new@example.com",
		Nickname: "新用户",
	}

	err := svc.Register(req)
	if err != nil {
		t.Errorf("Register() error = %v, want nil", err)
	}
}

func TestAuthServiceImpl_Register_UsernameExists(t *testing.T) {
	svc, mockRepo := setupTestService()

	// 预先添加用户
	hash, _ := bcrypt.HashPassword("password123")
	mockRepo.AddUser(1, "existinguser", "other@example.com", hash)

	req := request.RegisterRequest{
		Username: "existinguser",
		Password: "password123",
		Email:    "new@example.com",
	}

	err := svc.Register(req)
	if err == nil {
		t.Error("Register() 应该返回用户名已存在错误")
	}
	if err.Error() != "用户名已存在" {
		t.Errorf("error = %v, want '用户名已存在'", err)
	}
}

func TestAuthServiceImpl_Register_EmailExists(t *testing.T) {
	svc, mockRepo := setupTestService()

	// 预先添加用户
	hash, _ := bcrypt.HashPassword("password123")
	mockRepo.AddUser(1, "otheruser", "existing@example.com", hash)

	req := request.RegisterRequest{
		Username: "newuser",
		Password: "password123",
		Email:    "existing@example.com",
	}

	err := svc.Register(req)
	if err == nil {
		t.Error("Register() 应该返回邮箱已存在错误")
	}
	if err.Error() != "邮箱已存在" {
		t.Errorf("error = %v, want '邮箱已存在'", err)
	}
}

func TestAuthServiceImpl_Login_Success(t *testing.T) {
	svc, mockRepo := setupTestService()

	// 添加测试用户
	hash, _ := bcrypt.HashPassword("password123")
	mockRepo.AddUser(1, "testuser", "test@example.com", hash)

	req := request.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := svc.Login(req)
	if err != nil {
		t.Errorf("Login() error = %v, want nil", err)
	}
	if resp == nil {
		t.Fatal("Login() 返回 nil")
	}
	if resp.Token == "" {
		t.Error("token 不应为空")
	}
	if resp.User.Username != "testuser" {
		t.Errorf("username = %v, want 'testuser'", resp.User.Username)
	}
}

func TestAuthServiceImpl_Login_WrongPassword(t *testing.T) {
	svc, mockRepo := setupTestService()

	hash, _ := bcrypt.HashPassword("password123")
	mockRepo.AddUser(1, "testuser", "test@example.com", hash)

	req := request.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	_, err := svc.Login(req)
	if err == nil {
		t.Error("Login() 应该返回密码错误")
	}
	if err.Error() != "用户名或密码错误" {
		t.Errorf("error = %v, want '用户名或密码错误'", err)
	}
}

func TestAuthServiceImpl_Login_UserNotFound(t *testing.T) {
	svc, _ := setupTestService()

	req := request.LoginRequest{
		Username: "nonexistent",
		Password: "password123",
	}

	_, err := svc.Login(req)
	if err == nil {
		t.Error("Login() 应该返回用户不存在错误")
	}
}

func TestAuthServiceImpl_Login_UserDisabled(t *testing.T) {
	svc, mockRepo := setupTestService()

	hash, _ := bcrypt.HashPassword("password123")
	user := &entity.User{
		Username: "disableduser",
		Email:    "disabled@example.com",
		Password: hash,
		Status:   0, // 禁用状态
	}
	user.ID = 1
	mockRepo.Create(user)

	req := request.LoginRequest{
		Username: "disableduser",
		Password: "password123",
	}

	_, err := svc.Login(req)
	if err == nil {
		t.Error("Login() 应该返回用户已被禁用错误")
	}
	if err.Error() != "用户已被禁用" {
		t.Errorf("error = %v, want '用户已被禁用'", err)
	}
}

func TestAuthServiceImpl_RefreshToken_Success(t *testing.T) {
	svc, mockRepo := setupTestService()

	// 添加测试用户
	hash, _ := bcrypt.HashPassword("password123")
	mockRepo.AddUser(1, "testuser", "test@example.com", hash)

	// 生成有效 token
	oldToken, err := jwt.GenerateToken(1, "testuser")
	if err != nil {
		t.Fatalf("生成Token失败: %v", err)
	}

	newToken, err := svc.RefreshToken(oldToken)
	if err != nil {
		t.Errorf("RefreshToken() error = %v, want nil", err)
	}
	if newToken == "" {
		t.Error("newToken 不应为空")
	}

	// 验证新 token 可以解析
	claims, err := jwt.ParseToken(newToken)
	if err != nil {
		t.Errorf("新 token 无法解析: %v", err)
	}
	if claims.UserID != 1 {
		t.Errorf("claims.UserID = %v, want 1", claims.UserID)
	}
	if claims.Username != "testuser" {
		t.Errorf("claims.Username = %v, want 'testuser'", claims.Username)
	}
}

func TestAuthServiceImpl_RefreshToken_InvalidToken(t *testing.T) {
	svc, _ := setupTestService()

	_, err := svc.RefreshToken("invalid-token")
	if err == nil {
		t.Error("RefreshToken() 应该返回无效令牌错误")
	}
	if err.Error() != "无效的令牌" {
		t.Errorf("error = %v, want '无效的令牌'", err)
	}
}

func TestAuthServiceImpl_RefreshToken_UserNotFound(t *testing.T) {
	svc, _ := setupTestService()

	// 生成一个有效 token，但用户不存在
	token, err := jwt.GenerateToken(999, "nonexistent")
	if err != nil {
		t.Fatalf("生成Token失败: %v", err)
	}

	_, err = svc.RefreshToken(token)
	if err == nil {
		t.Error("RefreshToken() 应该返回用户不存在错误")
	}
	if err.Error() != "用户不存在" {
		t.Errorf("error = %v, want '用户不存在'", err)
	}
}

func TestAuthServiceImpl_RefreshToken_UserDisabled(t *testing.T) {
	svc, mockRepo := setupTestService()

	// 添加禁用用户
	hash, _ := bcrypt.HashPassword("password123")
	user := &entity.User{
		Username: "disableduser",
		Email:    "disabled@example.com",
		Password: hash,
		Status:   0,
	}
	user.ID = 1
	mockRepo.Create(user)

	token, err := jwt.GenerateToken(1, "disableduser")
	if err != nil {
		t.Fatalf("生成Token失败: %v", err)
	}

	_, err = svc.RefreshToken(token)
	if err == nil {
		t.Error("RefreshToken() 应该返回用户已被禁用错误")
	}
	if err.Error() != "用户已被禁用" {
		t.Errorf("error = %v, want '用户已被禁用'", err)
	}
}
