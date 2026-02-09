package dao

import (
	"errors"
	"test/models"

	"gorm.io/gorm"
)

func CreateExam(exam *models.Exam) error {
	return DB.Create(exam).Error
}

func CreateExamReviews(reviews []*models.ExamReview) error {
	return DB.Create(&reviews).Error
}

func GetPendingReviewsByReviewer(reviewerID uint) ([]models.ExamReview, error) {
	var reviews []models.ExamReview
	err := DB.Where("reviewer_id = ? AND status = ?", reviewerID, "pending").Find(&reviews).Error
	return reviews, err
}

func UpdateReviewWithVersion(review *models.ExamReview) error {
	result := DB.Model(&models.ExamReview{}).
		Where("id = ? AND version = ?", review.ID, review.Version).
		Updates(map[string]interface{}{
			"score":       review.Score,
			"comment":     review.Comment,
			"status":      review.Status,
			"reviewed_at": review.ReviewedAt,
			"version":     review.Version + 1,
		})

	if result.RowsAffected == 0 {
		return errors.New("记录已被修改，请刷新后重试")
	}
	return result.Error
}

func GetReviewsByExamAndStudent(examID, studentID uint) ([]models.ExamReview, error) {
	var reviews []models.ExamReview
	err := DB.Where("exam_id = ? AND student_id = ?", examID, studentID).Find(&reviews).Error
	return reviews, err
}

func GetReviewByID(id uint) (*models.ExamReview, error) {
	var review models.ExamReview
	err := DB.First(&review, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("阅卷记录不存在")
	}
	if err != nil {
		return nil, err
	}
	return &review, nil
}
