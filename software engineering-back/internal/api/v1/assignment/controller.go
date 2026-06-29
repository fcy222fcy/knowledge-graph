package assignment

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/service"
	"software_engineering/pkg/response"
)

// ─── 学生端 ────────────────────────────────────────────

// ListAssignments 学生查看作业列表
func ListAssignments(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	list, total, err := service.ListStudentAssignments(userID, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Paginated(c, list, total, page, size)
}

// GetAssignmentDetail 学生查看作业详情
func GetAssignmentDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	detail, err := service.GetStudentAssignmentDetail(uint(id))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, detail)
}

// SubmitAssignment 学生提交作业
func SubmitAssignment(c *gin.Context) {
	userID := c.GetUint("user_id")
	assignmentID, _ := strconv.Atoi(c.Param("id"))

	var req request.SubmitAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := service.SubmitAssignment(userID, uint(assignmentID), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, result)
}

// GetAssignmentResult 学生查看作业结果
func GetAssignmentResult(c *gin.Context) {
	userID := c.GetUint("user_id")
	assignmentID, _ := strconv.Atoi(c.Param("id"))

	result, err := service.GetStudentAssignmentResult(userID, uint(assignmentID))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, result)
}

// ─── 教师端（Admin） ──────────────────────────────────

// AdminListAssignments 教师查看作业列表
func AdminListAssignments(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	list, total, err := service.ListTeacherAssignments(userID, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Paginated(c, list, total, page, size)
}

// AdminCreateAssignment 教师创建作业
func AdminCreateAssignment(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := service.CreateAssignment(userID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, result)
}

// AdminUpdateAssignment 教师更新作业
func AdminUpdateAssignment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req request.UpdateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := service.UpdateAssignment(uint(id), req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, nil)
}

// AdminDeleteAssignment 教师删除作业
func AdminDeleteAssignment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.DeleteAssignment(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, nil)
}

// AdminPublishAssignment 发布作业
func AdminPublishAssignment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.PublishAssignment(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, nil)
}

// AdminCloseAssignment 关闭作业
func AdminCloseAssignment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := service.CloseAssignment(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, nil)
}

// AdminGetAssignment 教师查看作业详情
func AdminGetAssignment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	detail, err := service.GetAssignmentDetail(uint(id))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, detail)
}

// AdminListSubmissions 查看提交列表
func AdminListSubmissions(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	list, total, err := service.ListAssignmentSubmissions(uint(id), page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Paginated(c, list, total, page, size)
}
