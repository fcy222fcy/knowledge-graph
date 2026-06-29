package repository

import (
	"software_engineering/internal/model/entity"
	"software_engineering/pkg/database"
)

// ─── Assignment ────────────────────────────────────────

func CreateAssignment(a *entity.Assignment) error {
	return database.DB.Create(a).Error
}

func UpdateAssignment(a *entity.Assignment) error {
	return database.DB.Save(a).Error
}

func DeleteAssignment(id uint) error {
	return database.DB.Delete(&entity.Assignment{}, id).Error
}

func FindAssignmentByID(id uint) (*entity.Assignment, error) {
	var a entity.Assignment
	err := database.DB.First(&a, id).Error
	return &a, err
}

func ListAssignmentsByTeacher(teacherID uint, page, size int) ([]entity.Assignment, int64, error) {
	var list []entity.Assignment
	var total int64
	query := database.DB.Where("teacher_id = ?", teacherID)
	query.Model(&entity.Assignment{}).Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&list).Error
	return list, total, err
}

func ListPublishedAssignments(page, size int) ([]entity.Assignment, int64, error) {
	var list []entity.Assignment
	var total int64
	query := database.DB.Where("status = ?", "published")
	query.Model(&entity.Assignment{}).Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&list).Error
	return list, total, err
}

// ─── AssignmentQuestion ────────────────────────────────

func CreateAssignmentQuestions(questions []entity.AssignmentQuestion) error {
	if len(questions) == 0 {
		return nil
	}
	return database.DB.Create(&questions).Error
}

func DeleteAssignmentQuestions(assignmentID uint) error {
	return database.DB.Where("assignment_id = ?", assignmentID).Delete(&entity.AssignmentQuestion{}).Error
}

func ListAssignmentQuestions(assignmentID uint) ([]entity.AssignmentQuestion, error) {
	var list []entity.AssignmentQuestion
	err := database.DB.Where("assignment_id = ?", assignmentID).Order("sort_order ASC, id ASC").Find(&list).Error
	return list, err
}

func FindAssignmentQuestionByID(id uint) (*entity.AssignmentQuestion, error) {
	var q entity.AssignmentQuestion
	err := database.DB.First(&q, id).Error
	return &q, err
}

// ─── AssignmentSubmission ──────────────────────────────

func CreateAssignmentSubmission(s *entity.AssignmentSubmission) error {
	return database.DB.Create(s).Error
}

func FindAssignmentSubmission(assignmentID, userID uint) (*entity.AssignmentSubmission, error) {
	var s entity.AssignmentSubmission
	err := database.DB.Where("assignment_id = ? AND user_id = ?", assignmentID, userID).First(&s).Error
	return &s, err
}

func ListAssignmentSubmissions(assignmentID uint, page, size int) ([]entity.AssignmentSubmission, int64, error) {
	var list []entity.AssignmentSubmission
	var total int64
	query := database.DB.Where("assignment_id = ?", assignmentID)
	query.Model(&entity.AssignmentSubmission{}).Count(&total)
	err := query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&list).Error
	return list, total, err
}

func CountAssignmentSubmissions(assignmentID uint) int {
	var count int64
	database.DB.Model(&entity.AssignmentSubmission{}).Where("assignment_id = ?", assignmentID).Count(&count)
	return int(count)
}

func CountAssignmentQuestions(assignmentID uint) int {
	var count int64
	database.DB.Model(&entity.AssignmentQuestion{}).Where("assignment_id = ?", assignmentID).Count(&count)
	return int(count)
}

// CountAssignments 统计作业总数
func CountAssignments() (int64, error) {
	var count int64
	err := database.DB.Model(&entity.Assignment{}).Count(&count).Error
	return count, err
}
