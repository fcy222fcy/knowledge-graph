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

// TeacherRegister 教师注册，验证用户名和邮箱唯一性后创建教师账号
func TeacherRegister(req request.TeacherRegisterRequest) error {
	existing, _ := repository.FindTeacherByUsername(req.Username)
	if existing.ID != 0 {
		return errors.New("用户名已存在")
	}
	existingEmail, _ := repository.FindTeacherByEmail(req.Email)
	if existingEmail.ID != 0 {
		return errors.New("邮箱已存在")
	}

	hash, err := bcrypt.HashPassword(req.Password)
	if err != nil {
		return err
	}

	teacher := &entity.Teacher{
		Username: req.Username,
		Password: hash,
		Email:    req.Email,
		Nickname: req.Nickname,
		Status:   1,
	}
	return repository.CreateTeacher(teacher)
}

// TeacherLogin 教师登录，验证用户名密码后生成 JWT Token
func TeacherLogin(req request.TeacherLoginRequest) (*response.TeacherLoginResponse, error) {
	teacher, err := repository.FindTeacherByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	if teacher.Status == 0 {
		return nil, errors.New("账号已被禁用")
	}

	if !bcrypt.CheckPassword(req.Password, teacher.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	token, err := jwt.GenerateToken(teacher.ID, teacher.Username, "teacher")
	if err != nil {
		return nil, err
	}

	return &response.TeacherLoginResponse{
		Token: token,
		Teacher: response.TeacherResponse{
			ID:        teacher.ID,
			Username:  teacher.Username,
			Email:     teacher.Email,
			Nickname:  teacher.Nickname,
			Avatar:    teacher.Avatar,
			Status:    teacher.Status,
			CreatedAt: teacher.CreatedAt.Format("2006-01-02T15:04:05Z"),
		},
	}, nil
}

// TeacherRefreshToken 刷新教师 JWT Token，验证旧 Token 有效性后生成新 Token
func TeacherRefreshToken(oldToken string) (string, error) {
	claims, err := jwt.ParseToken(oldToken)
	if err != nil {
		return "", errors.New("无效的令牌")
	}
	teacher, err := repository.FindTeacherByID(claims.UserID)
	if err != nil {
		return "", errors.New("教师不存在")
	}
	if teacher.Status == 0 {
		return "", errors.New("账号已被禁用")
	}
	return jwt.GenerateToken(teacher.ID, teacher.Username, "teacher")
}
