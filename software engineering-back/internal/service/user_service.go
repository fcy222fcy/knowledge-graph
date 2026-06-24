package service

import (
	"errors"

	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/model/dto/response"
	"software_engineering/internal/repository"
	"software_engineering/pkg/bcrypt"
)

// GetProfile 获取用户个人信息
func GetProfile(userID uint) (*response.UserResponse, error) {
	user, err := repository.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return &response.UserResponse{
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

// UpdateProfile 更新用户个人信息，仅更新非空字段
func UpdateProfile(userID uint, req request.UpdateProfileRequest) error {
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

// ChangePassword 修改用户密码，需验证旧密码
func ChangePassword(userID uint, req request.ChangePasswordRequest) error {
	user, err := repository.FindUserByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}
	if !bcrypt.CheckPassword(req.OldPassword, user.Password) {
		return errors.New("旧密码错误")
	}
	hash, err := bcrypt.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	user.Password = hash
	return repository.UpdateUser(user)
}

// ListUsers 分页获取用户列表
func ListUsers(page, size int) (*response.UserListResponse, error) {
	users, total, err := repository.ListUsers(page, size)
	if err != nil {
		return nil, err
	}
	list := make([]response.UserResponse, len(users))
	for i, u := range users {
		list[i] = response.UserResponse{
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
	return &response.UserListResponse{List: list, Total: total, Page: page, Size: size, TotalPage: totalPage}, nil
}
