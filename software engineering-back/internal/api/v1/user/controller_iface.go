package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/model/dto/response"
	pkgResponse "software_engineering/pkg/response"
)

// UserService 定义用户服务接口
type UserService interface {
	GetProfile(userID uint) (*response.UserResponse, error)                                 // 获取用户资料
	UpdateProfile(userID uint, req request.UpdateProfileRequest) error                      // 更新用户资料
	ChangePassword(userID uint, req request.ChangePasswordRequest) error                    // 修改密码
	ListUsers(page, size int) (*response.UserListResponse, error)                          // 列表用户
}

// UserController 用户控制器
type UserController struct {
	userService UserService // 用户服务
}

// NewUserController 创建用户控制器实例
func NewUserController(userService UserService) *UserController {
	return &UserController{userService: userService}
}

// GetProfile 获取用户个人信息
func (ctrl *UserController) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	resp, err := ctrl.userService.GetProfile(userID)
	if err != nil {
		pkgResponse.Error(c, http.StatusNotFound, err.Error())
		return
	}
	pkgResponse.Success(c, resp)
}

// UpdateProfile 更新用户个人信息
func (ctrl *UserController) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := ctrl.userService.UpdateProfile(userID, req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	pkgResponse.Success(c, nil)
}

// ChangePassword 修改用户密码
func (ctrl *UserController) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	if err := ctrl.userService.ChangePassword(userID, req); err != nil {
		pkgResponse.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	pkgResponse.Success(c, nil)
}

// ListUsers 获取用户列表
func (ctrl *UserController) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	resp, err := ctrl.userService.ListUsers(page, size)
	if err != nil {
		pkgResponse.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkgResponse.Success(c, resp)
}
