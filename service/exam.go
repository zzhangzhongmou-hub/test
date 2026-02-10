package service

import (
	"errors"
	"fmt"
	"test/dao"
	"test/models"
	"time"
)

type CreateExamRequest struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	Department    string `json:"department"`
	Deadline      string `json:"deadline"`
	ReviewerCount int    `json:"reviewer_count"`
	Duration      int    `json:"duration"`
	CreatorID     uint
}

func CreateExam(req CreateExamRequest) error {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	deadline, err := time.ParseInLocation("2006-01-02 15:04:05", req.Deadline, loc)
	if err != nil {
		return errors.New("截止时间格式错误，请使用 2006-01-02 15:04:05 格式")
	}

	if deadline.Before(time.Now()) {
		return errors.New("截止时间必须晚于当前时间")
	}

	exam := &models.Exam{
		Title:       req.Title,
		Description: req.Description,
		Department:  req.Department,
		CreatorID:   req.CreatorID,
		Deadline:    deadline,
		TotalScore:  100,
		Duration:    req.Duration,
		Status:      "pending",
	}

	if err := dao.CreateExam(exam); err != nil {
		return err
	}

	return assignReviewers(exam.ID, req.Department, req.ReviewerCount)
}

func assignReviewers(examID uint, dept string, reviewerCount int) error {
	admins, err := dao.GetUsersByDepartmentAndRole(dept, "admin")
	if err != nil || len(admins) == 0 {
		return errors.New("该部门无可用阅卷人")
	}

	students, err := dao.GetUsersByDepartmentAndRole(dept, "student")
	if err != nil {
		return err
	}

	var reviews []*models.ExamReview
	for i, student := range students {
		for j := 0; j < reviewerCount; j++ {
			adminIndex := (i + j) % len(admins)
			reviews = append(reviews, &models.ExamReview{
				ExamID:     examID,
				StudentID:  student.ID,
				ReviewerID: admins[adminIndex].ID,
				Status:     "pending",
				AssignedAt: time.Now(),
			})
		}
	}

	return dao.CreateExamReviews(reviews)
}

type SubmitReviewRequest struct {
	ReviewID   uint
	ReviewerID uint
	Score      int
	Comment    string
	Version    int
}

func SubmitReview(req SubmitReviewRequest) error {
	review, err := dao.GetReviewByID(req.ReviewID)
	if err != nil {
		return err
	}
	if review.ReviewerID != req.ReviewerID {
		return errors.New("无权批改此份考卷")
	}
	if review.Version != req.Version {
		return errors.New("考卷已被其他管理员修改，请刷新后重试")
	}
	if review.Status != "pending" {
		return errors.New("该考卷已被批改")
	}
	now := time.Now()
	review.Score = &req.Score
	review.Comment = req.Comment
	review.Status = "completed"
	review.ReviewedAt = &now
	if err := dao.UpdateReviewWithVersion(review); err != nil {
		return err
	}
	go calculateFinalScore(review.ExamID, review.StudentID)

	return nil
}

func calculateFinalScore(examID, studentID uint) {
	reviews, err := dao.GetReviewsByExamAndStudent(examID, studentID)
	if err != nil {
		return
	}
	totalScore := 0
	count := 0
	for _, r := range reviews {
		if r.Status != "completed" || r.Score == nil {
			return
		}
		totalScore += *r.Score
		count++
	}
	if count > 0 {
		finalScore := totalScore / count
		fmt.Printf("[Exam] 学生%d的最终成绩: %d（由%d位老师评分）\n",
			studentID, finalScore, count)
	}
}
