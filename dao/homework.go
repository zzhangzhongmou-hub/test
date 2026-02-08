package dao

import (
	"errors"
	"test/models"
	"time"

	"gorm.io/gorm"
)

func CreateHomework(homework *models.Homework) error {

	return DB.Create(homework).Error
}

func GetHomeworkByID(id uint) (*models.Homework, error) {
	var hw models.Homework
	err := DB.First(&hw, id).Error
	if err != nil {
		return nil, err
	}
	return &hw, nil
}

func GetHomeworksByDepartment(department string, page, pageSize int) ([]models.Homework, int64, error) {
	var homeworks []models.Homework
	var total int64

	query := DB.Model(&models.Homework{})

	if department != "" {
		query = query.Where("department = ?", department)
	}
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&homeworks).Error

	return homeworks, total, err
}

func GetHomeworksByDeadlineRange(start, end time.Time) ([]models.Homework, error) {
	var homeworks []models.Homework
	err := DB.Where("deadline BETWEEN ? AND ?", start, end).Find(&homeworks).Error
	return homeworks, err
}

func UpdateHomework(homework *models.Homework) error {
	result := DB.Model(&models.Homework{}).
		Where("id = ? AND version = ?", homework.ID, homework.Version).
		Updates(map[string]interface{}{
			"title":       homework.Title,
			"description": homework.Description,
			"deadline":    homework.Deadline,
			"allow_late":  homework.AllowLate,
			"version":     gorm.Expr("version + 1"),
		})
	if result.RowsAffected == 0 {
		return errors.New("作业已被其他管理员修改，请刷新后重试")
	}
	return result.Error
}
func DeleteHomework(id uint) error {
	return DB.Delete(&models.Homework{}, id).Error
}
