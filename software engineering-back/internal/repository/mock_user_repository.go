package repository

import (
	"software_engineering/internal/model/entity"

	"gorm.io/gorm"
)

// MockUserRepository 用户仓库的 Mock 实现
type MockUserRepository struct {
	Users   map[uint]*entity.User
	NextID  uint
	Err     error // 模拟错误
}

// NewMockUserRepository 创建 Mock 用户仓库
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Users:  make(map[uint]*entity.User),
		NextID: 1,
	}
}

func (m *MockUserRepository) Create(user *entity.User) error {
	if m.Err != nil {
		return m.Err
	}
	user.ID = m.NextID
	m.NextID++
	m.Users[user.ID] = user
	return nil
}

func (m *MockUserRepository) FindByUsername(username string) (*entity.User, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	for _, user := range m.Users {
		if user.Username == username {
			return user, nil
		}
	}
	// 返回空用户而非 nil，与 GORM 行为一致
	return &entity.User{}, gorm.ErrRecordNotFound
}

func (m *MockUserRepository) FindByEmail(email string) (*entity.User, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	for _, user := range m.Users {
		if user.Email == email {
			return user, nil
		}
	}
	// 返回空用户而非 nil，与 GORM 行为一致
	return &entity.User{}, gorm.ErrRecordNotFound
}

func (m *MockUserRepository) FindByID(id uint) (*entity.User, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	user, ok := m.Users[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (m *MockUserRepository) Update(user *entity.User) error {
	if m.Err != nil {
		return m.Err
	}
	m.Users[user.ID] = user
	return nil
}

func (m *MockUserRepository) List(page, size int) ([]entity.User, int64, error) {
	if m.Err != nil {
		return nil, 0, m.Err
	}
	var users []entity.User
	for _, user := range m.Users {
		users = append(users, *user)
	}
	return users, int64(len(users)), nil
}

// AddUser 辅助方法：添加测试用户
func (m *MockUserRepository) AddUser(id uint, username, email, password string) {
	user := &entity.User{
		Username: username,
		Email:    email,
		Password: password,
		Status:   1,
	}
	user.ID = id // 通过 BaseModel 设置 ID
	m.Users[id] = user
	if id >= m.NextID {
		m.NextID = id + 1
	}
}

// SetError 辅助方法：设置错误
func (m *MockUserRepository) SetError(err error) {
	m.Err = err
}

// ClearError 辅助方法：清除错误
func (m *MockUserRepository) ClearError() {
	m.Err = nil
}

// MockAuthService 认证服务的 Mock 实现
type MockAuthService struct {
	RegisterFunc  func(req interface{}) error
	LoginFunc     func(req interface{}) (interface{}, error)
	RefreshFunc   func(token string) (string, error)
	Err           error
}

func (m *MockAuthService) Register(req interface{}) error {
	if m.Err != nil {
		return m.Err
	}
	if m.RegisterFunc != nil {
		return m.RegisterFunc(req)
	}
	return nil
}

func (m *MockAuthService) Login(req interface{}) (interface{}, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if m.LoginFunc != nil {
		return m.LoginFunc(req)
	}
	return nil, nil
}

func (m *MockAuthService) RefreshToken(oldToken string) (string, error) {
	if m.Err != nil {
		return "", m.Err
	}
	if m.RefreshFunc != nil {
		return m.RefreshFunc(oldToken)
	}
	return "new-token", nil
}
