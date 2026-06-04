package service

import (
	"errors"

	"software_engineering/internal/model/dto"
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
