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

func Register(req request.RegisterRequest) error {
	existing, _ := repository.FindUserByUsername(req.Username)
	if existing.ID != 0 {
		return errors.New("用户名已存在")
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
	return repository.CreateUser(user)
}

func Login(req request.LoginRequest) (*response.LoginResponse, error) {
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

func RefreshToken(oldToken string) (string, error) {
	claims, err := jwt.ParseToken(oldToken)
	if err != nil {
		return "", errors.New("无效的令牌")
	}
	return jwt.GenerateToken(claims.UserID, claims.Username)
}
