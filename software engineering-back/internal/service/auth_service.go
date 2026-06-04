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
