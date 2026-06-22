package service

import (
	"errors"

	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/model/dto/response"
	"software_engineering/internal/model/entity"
	"software_engineering/internal/repository"
	"software_engineering/pkg/bcrypt"
	"software_engineering/pkg/jwt"

	"gorm.io/gorm"
)

// AuthService 定义认证服务接口
type AuthService interface {
	Register(req request.RegisterRequest) error                        // 用户注册
	Login(req request.LoginRequest) (*response.LoginResponse, error)   // 用户登录
	RefreshToken(oldToken string) (string, error)                     // 刷新令牌
}

// AuthServiceImpl 认证服务实现
type AuthServiceImpl struct {
	userRepo repository.UserRepository // 用户仓库
}

// NewAuthService 创建认证服务实例
func NewAuthService(userRepo repository.UserRepository) *AuthServiceImpl {
	return &AuthServiceImpl{userRepo: userRepo}
}

// Register 用户注册，验证用户名和邮箱唯一性后创建用户
func (s *AuthServiceImpl) Register(req request.RegisterRequest) error {
	existing, _ := s.userRepo.FindByUsername(req.Username)
	if existing.ID != 0 {
		return errors.New("用户名已存在")
	}
	existingEmail, _ := s.userRepo.FindByEmail(req.Email)
	if existingEmail.ID != 0 {
		return errors.New("邮箱已存在")
	}

	hash, err := bcrypt.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &entity.User{
		Username: req.Username,
		Password: hash,
		Email:    req.Email,
		Nickname: req.Nickname,
		Status:   1,
	}
	return s.userRepo.Create(user)
}

// Login 用户登录，验证用户名密码后生成 JWT Token
func (s *AuthServiceImpl) Login(req request.LoginRequest) (*response.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	if user.Status == 0 {
		return nil, errors.New("用户已被禁用")
	}

	if !bcrypt.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &response.LoginResponse{
		Token: token,
		User: response.UserResponse{
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

// RefreshToken 刷新 JWT Token，验证旧 Token 有效性后生成新 Token
func (s *AuthServiceImpl) RefreshToken(oldToken string) (string, error) {
	claims, err := jwt.ParseToken(oldToken)
	if err != nil {
		return "", errors.New("无效的令牌")
	}
	// 验证用户是否存在且未被禁用
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return "", errors.New("用户不存在")
	}
	if user.Status == 0 {
		return "", errors.New("用户已被禁用")
	}
	return jwt.GenerateToken(claims.UserID, claims.Username)
}
