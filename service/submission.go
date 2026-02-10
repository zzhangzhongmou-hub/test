package service

import (
	"errors"
	"test/dao"
	"test/models"
	"time"
)

type SubmitRequest struct {
	HomeworkID uint `json:"homework_id"`
	StudentID  uint
	Content    string `json:"content" binding:"required"`
	FileURL    string `json:"file_url"`
}

func Submit(req SubmitRequest) error {
	homework, err := dao.GetHomeworkByID(req.HomeworkID)
	if err != nil {
		return errors.New("作业不存在")
	}

	student, _ := dao.GetUserByID(req.StudentID)
	if student.Department != homework.Department {
		return errors.New("只能提交本部门的作业")
	}

	exists, _ := dao.CheckSubmissionExists(req.HomeworkID, req.StudentID)
	if exists {
		return errors.New("已提交过该作业，请勿重复提交")
	}

	isLate := time.Now().After(homework.Deadline)
	if isLate && !homework.AllowLate {
		return errors.New("已过截止时间，且该作业不允许补交")
	}

	submission := &models.Submission{
		HomeworkID: req.HomeworkID,
		StudentID:  req.StudentID,
		Content:    req.Content,
		FileURL:    req.FileURL,
		IsLate:     isLate,
		Status:     "submitted",
	}

	return dao.CreateSubmission(submission)
}

func GetMySubmissions(studentID uint, page, pageSize int) ([]models.Submission, int64, error) {
	return dao.GetSubmissionsByStudent(studentID, page, pageSize)
}

func GetSubmissionsByHomework(homeworkID uint, page, pageSize int) ([]models.Submission, int64, error) {
	return dao.GetSubmissionsByHomework(homeworkID, page, pageSize)
}

type ReviewRequest struct {
	SubmissionID uint `json:"submission_id"`
	ReviewerID   uint
	ReviewerDept string
	Score        int    `json:"score" binding:"required,min=0,max=100"`
	Comment      string `json:"comment" binding:"required"`
}

func Review(req ReviewRequest) error {
	sub, err := dao.GetSubmissionByID(req.SubmissionID)
	if err != nil {
		return errors.New("提交记录不存在")
	}

	homework, _ := dao.GetHomeworkByID(sub.HomeworkID)

	if homework.Department != req.ReviewerDept {
		return errors.New("只能批改本部门的作业")
	}

	sub.Score = &req.Score
	sub.Comment = req.Comment
	sub.Status = "graded"
	now := time.Now()
	sub.ReviewedAt = &now

	if err := dao.UpdateSubmission(sub); err != nil {
		return err
	}
	go func() {
		student, _ := dao.GetUserByID(sub.StudentID)
		if student.Email != "" {
			EmailSvc.NotifyReview(student.ID, homework.Title, req.Score, req.Comment)
		}
	}()

	return nil
}

type MarkExcellentRequest struct {
	SubmissionID uint
	IsExcellent  bool
	ReviewerDept string
}

func MarkExcellent(req MarkExcellentRequest) error {
	sub, err := dao.GetSubmissionByID(req.SubmissionID)
	if err != nil {
		return err
	}

	homework, _ := dao.GetHomeworkByID(sub.HomeworkID)
	if homework.Department != req.ReviewerDept {
		return errors.New("只能标记本部门的优秀作业")
	}

	return dao.MarkSubmissionExcellent(req.SubmissionID, req.IsExcellent)
}

func GetExcellentSubmissions(department string, page, pageSize int) ([]models.Submission, int64, error) {
	return dao.GetExcellentSubmissions(department, page, pageSize)
}
