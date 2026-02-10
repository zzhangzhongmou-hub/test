package dao

import (
	"errors"
	"test/models"

	"gorm.io/gorm"
)

func CreateSubmission(submission *models.Submission) error {
	return DB.Create(submission).Error
}

func GetSubmissionByID(id uint) (*models.Submission, error) {
	var sub models.Submission
	err := DB.First(&sub, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("提交记录不存在")
	}
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func GetSubmissionsByStudent(studentID uint, page, pageSize int) ([]models.Submission, int64, error) {
	var submissions []models.Submission
	var total int64

	query := DB.Model(&models.Submission{}).Where("student_id = ?", studentID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&submissions).Error

	return submissions, total, err
}

func GetSubmissionsByHomework(homeworkID uint, page, pageSize int) ([]models.Submission, int64, error) {
	var submissions []models.Submission
	var total int64

	query := DB.Model(&models.Submission{}).Where("homework_id = ?", homeworkID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&submissions).Error

	return submissions, total, err
}

func UpdateSubmission(submission *models.Submission) error {
	return DB.Save(submission).Error
}

func CheckSubmissionExists(homeworkID, studentID uint) (bool, error) {
	var count int64
	err := DB.Model(&models.Submission{}).
		Where("homework_id = ? AND student_id = ?", homeworkID, studentID).
		Count(&count).Error
	return count > 0, err
}

func MarkSubmissionExcellent(submissionID uint, isExcellent bool) error {
	return DB.Model(&models.Submission{}).
		Where("id = ?", submissionID).
		Update("is_excellent", isExcellent).Error
}

func GetExcellentSubmissions(department string, page, pageSize int) ([]models.Submission, int64, error) {
	var submissions []models.Submission
	var total int64

	query := DB.Model(&models.Submission{}).
		Where("is_excellent = ?", true).
		Preload("Homework").Preload("Student")

	if department != "" {
		query = query.Joins("JOIN homeworks ON submissions.homework_id = homeworks.id").
			Where("homeworks.department = ?", department)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("submissions.created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&submissions).Error

	return submissions, total, err
}
