package service

import (
	"errors"
	"testing"

	"software_engineering/internal/model/dto/request"
	"software_engineering/internal/repository"
)

func setupQuizTestService() (*QuizServiceImpl, *repository.MockQuestionRepository, *repository.MockQuizRepository) {
	mockQuestionRepo := repository.NewMockQuestionRepository()
	mockQuizRepo := repository.NewMockQuizRepository()
	svc := NewQuizService(mockQuestionRepo, mockQuizRepo)
	return svc, mockQuestionRepo, mockQuizRepo
}

func TestQuizServiceImpl_SubmitQuiz_Success(t *testing.T) {
	svc, mockQuestionRepo, _ := setupQuizTestService()

	// 添加测试题目
	mockQuestionRepo.AddQuestion(1, "1+1=?", "A", "single", "easy",
		`[{"key":"A","value":"2"},{"key":"B","value":"3"}]`, "基础加法")

	req := request.SubmitQuizRequest{
		QuestionID: 1,
		UserAnswer: "A",
	}

	resp, err := svc.SubmitQuiz(1, req)
	if err != nil {
		t.Errorf("SubmitQuiz() error = %v, want nil", err)
	}
	if resp == nil {
		t.Fatal("SubmitQuiz() 返回 nil")
	}
	if !resp.IsCorrect {
		t.Error("IsCorrect = false, want true")
	}
	if resp.CorrectAnswer != "A" {
		t.Errorf("CorrectAnswer = %v, want 'A'", resp.CorrectAnswer)
	}
}

func TestQuizServiceImpl_SubmitQuiz_WrongAnswer(t *testing.T) {
	svc, mockQuestionRepo, _ := setupQuizTestService()

	// 添加测试题目
	mockQuestionRepo.AddQuestion(1, "1+1=?", "A", "single", "easy",
		`[{"key":"A","value":"2"},{"key":"B","value":"3"}]`, "基础加法")

	req := request.SubmitQuizRequest{
		QuestionID: 1,
		UserAnswer: "B",
	}

	resp, err := svc.SubmitQuiz(1, req)
	if err != nil {
		t.Errorf("SubmitQuiz() error = %v, want nil", err)
	}
	if resp == nil {
		t.Fatal("SubmitQuiz() 返回 nil")
	}
	if resp.IsCorrect {
		t.Error("IsCorrect = true, want false")
	}
}

func TestQuizServiceImpl_SubmitQuiz_QuestionNotFound(t *testing.T) {
	svc, _, _ := setupQuizTestService()

	req := request.SubmitQuizRequest{
		QuestionID: 999,
		UserAnswer: "A",
	}

	_, err := svc.SubmitQuiz(1, req)
	if err == nil {
		t.Error("SubmitQuiz() 应该返回题目不存在错误")
	}
	if err.Error() != "题目不存在" {
		t.Errorf("error = %v, want '题目不存在'", err)
	}
}

func TestQuizServiceImpl_SubmitQuiz_CreateQuizError(t *testing.T) {
	svc, mockQuestionRepo, mockQuizRepo := setupQuizTestService()

	// 添加测试题目
	mockQuestionRepo.AddQuestion(1, "1+1=?", "A", "single", "easy",
		`[{"key":"A","value":"2"},{"key":"B","value":"3"}]`, "基础加法")

	// 设置 mock 错误
	mockQuizRepo.SetError(errors.New("数据库错误"))

	req := request.SubmitQuizRequest{
		QuestionID: 1,
		UserAnswer: "A",
	}

	_, err := svc.SubmitQuiz(1, req)
	if err == nil {
		t.Error("SubmitQuiz() 应该返回数据库错误")
	}
	if err.Error() != "数据库错误" {
		t.Errorf("error = %v, want '数据库错误'", err)
	}
}

func TestQuizServiceImpl_GetQuizDetail_Success(t *testing.T) {
	svc, mockQuestionRepo, mockQuizRepo := setupQuizTestService()

	// 添加测试题目
	mockQuestionRepo.AddQuestion(1, "1+1=?", "A", "single", "easy",
		`[{"key":"A","value":"2"},{"key":"B","value":"3"}]`, "基础加法")

	// 添加测试答题记录
	mockQuizRepo.AddQuiz(1, 1, 1, "A", true)

	resp, err := svc.GetQuizDetail(1)
	if err != nil {
		t.Errorf("GetQuizDetail() error = %v, want nil", err)
	}
	if resp == nil {
		t.Fatal("GetQuizDetail() 返回 nil")
	}
	if resp.QuestionTitle != "1+1=?" {
		t.Errorf("QuestionTitle = %v, want '1+1=?'", resp.QuestionTitle)
	}
	if resp.UserAnswer != "A" {
		t.Errorf("UserAnswer = %v, want 'A'", resp.UserAnswer)
	}
	if !resp.IsCorrect {
		t.Error("IsCorrect = false, want true")
	}
}

func TestQuizServiceImpl_GetQuizDetail_QuizNotFound(t *testing.T) {
	svc, _, _ := setupQuizTestService()

	_, err := svc.GetQuizDetail(999)
	if err == nil {
		t.Error("GetQuizDetail() 应该返回答题记录不存在错误")
	}
	if err.Error() != "答题记录不存在" {
		t.Errorf("error = %v, want '答题记录不存在'", err)
	}
}

func TestQuizServiceImpl_GetQuizDetail_QuestionNotFound(t *testing.T) {
	svc, _, mockQuizRepo := setupQuizTestService()

	// 添加答题记录，但题目不存在
	mockQuizRepo.AddQuiz(1, 999, 1, "A", true)

	_, err := svc.GetQuizDetail(1)
	if err == nil {
		t.Error("GetQuizDetail() 应该返回题目不存在错误")
	}
	if err.Error() != "题目不存在" {
		t.Errorf("error = %v, want '题目不存在'", err)
	}
}

func TestQuizServiceImpl_ListQuizHistory_Success(t *testing.T) {
	svc, mockQuestionRepo, mockQuizRepo := setupQuizTestService()

	// 添加测试题目
	mockQuestionRepo.AddQuestion(1, "1+1=?", "A", "single", "easy",
		`[{"key":"A","value":"2"},{"key":"B","value":"3"}]`, "基础加法")
	mockQuestionRepo.AddQuestion(2, "2+2=?", "B", "single", "easy",
		`[{"key":"A","value":"3"},{"key":"B","value":"4"}]`, "基础加法")

	// 添加测试答题记录
	mockQuizRepo.AddQuiz(1, 1, 1, "A", true)
	mockQuizRepo.AddQuiz(2, 2, 1, "B", true)
	mockQuizRepo.AddQuiz(3, 1, 1, "B", false)

	list, total, err := svc.ListQuizHistory(1, 1, 10, 0, nil)
	if err != nil {
		t.Errorf("ListQuizHistory() error = %v, want nil", err)
	}
	if total != 3 {
		t.Errorf("total = %v, want 3", total)
	}
	if len(list) != 3 {
		t.Errorf("list length = %v, want 3", len(list))
	}
}

func TestQuizServiceImpl_ListQuizHistory_FilterByCorrect(t *testing.T) {
	svc, mockQuestionRepo, mockQuizRepo := setupQuizTestService()

	// 添加测试题目
	mockQuestionRepo.AddQuestion(1, "1+1=?", "A", "single", "easy",
		`[{"key":"A","value":"2"},{"key":"B","value":"3"}]`, "基础加法")

	// 添加测试答题记录
	mockQuizRepo.AddQuiz(1, 1, 1, "A", true)
	mockQuizRepo.AddQuiz(2, 1, 1, "B", false)

	// 只查看答对的记录
	isCorrect := true
	list, total, err := svc.ListQuizHistory(1, 1, 10, 0, &isCorrect)
	if err != nil {
		t.Errorf("ListQuizHistory() error = %v, want nil", err)
	}
	if total != 1 {
		t.Errorf("total = %v, want 1", total)
	}
	if len(list) != 1 {
		t.Errorf("list length = %v, want 1", len(list))
	}
	if !list[0].IsCorrect {
		t.Error("IsCorrect = false, want true")
	}
}

func TestQuizServiceImpl_ListQuizHistory_FilterByIncorrect(t *testing.T) {
	svc, mockQuestionRepo, mockQuizRepo := setupQuizTestService()

	// 添加测试题目
	mockQuestionRepo.AddQuestion(1, "1+1=?", "A", "single", "easy",
		`[{"key":"A","value":"2"},{"key":"B","value":"3"}]`, "基础加法")

	// 添加测试答题记录
	mockQuizRepo.AddQuiz(1, 1, 1, "A", true)
	mockQuizRepo.AddQuiz(2, 1, 1, "B", false)

	// 只查看答错的记录
	isCorrect := false
	list, total, err := svc.ListQuizHistory(1, 1, 10, 0, &isCorrect)
	if err != nil {
		t.Errorf("ListQuizHistory() error = %v, want nil", err)
	}
	if total != 1 {
		t.Errorf("total = %v, want 1", total)
	}
	if len(list) != 1 {
		t.Errorf("list length = %v, want 1", len(list))
	}
	if list[0].IsCorrect {
		t.Error("IsCorrect = true, want false")
	}
}

func TestQuizServiceImpl_ListQuizHistory_FilterByUser(t *testing.T) {
	svc, mockQuestionRepo, mockQuizRepo := setupQuizTestService()

	// 添加测试题目
	mockQuestionRepo.AddQuestion(1, "1+1=?", "A", "single", "easy",
		`[{"key":"A","value":"2"},{"key":"B","value":"3"}]`, "基础加法")

	// 添加不同用户的答题记录
	mockQuizRepo.AddQuiz(1, 1, 1, "A", true)
	mockQuizRepo.AddQuiz(2, 1, 2, "B", false)

	list, total, err := svc.ListQuizHistory(1, 1, 10, 0, nil)
	if err != nil {
		t.Errorf("ListQuizHistory() error = %v, want nil", err)
	}
	if total != 1 {
		t.Errorf("total = %v, want 1", total)
	}
	if len(list) != 1 {
		t.Errorf("list length = %v, want 1", len(list))
	}
}

func TestQuizServiceImpl_ListQuizHistory_EmptyResult(t *testing.T) {
	svc, _, _ := setupQuizTestService()

	list, total, err := svc.ListQuizHistory(999, 1, 10, 0, nil)
	if err != nil {
		t.Errorf("ListQuizHistory() error = %v, want nil", err)
	}
	if total != 0 {
		t.Errorf("total = %v, want 0", total)
	}
	if len(list) != 0 {
		t.Errorf("list length = %v, want 0", len(list))
	}
}

func TestQuizServiceImpl_ListQuizHistory_DatabaseError(t *testing.T) {
	svc, _, mockQuizRepo := setupQuizTestService()

	// 设置 mock 错误
	mockQuizRepo.SetError(errors.New("数据库错误"))

	_, _, err := svc.ListQuizHistory(1, 1, 10, 0, nil)
	if err == nil {
		t.Error("ListQuizHistory() 应该返回数据库错误")
	}
	if err.Error() != "数据库错误" {
		t.Errorf("error = %v, want '数据库错误'", err)
	}
}
