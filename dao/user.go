package dao

import (
	"test/models"

	"gorm.io/gorm"
)

func CreateUser(user *models.User) error {
	return DB.Create(user).Error
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUnsubmittedStudents(homeworkID uint, department string) ([]models.User, error) {
	var students []models.User

	subQuery := DB.Model(&models.Submission{}).
		Select("student_id").
		Where("homework_id = ?", homeworkID)
	
	err := DB.Where("department = ? AND role = ?", department, "student").
		Where("id NOT IN (?)", subQuery).
		Find(&students).Error

	return students, err
}

func UpdateUser(user *models.User) error {
	return DB.Save(user).Error
}

func DeleteUser(id uint) error {
	return DB.Delete(&models.User{}, id).Error
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := DB.Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUserEmail(userID uint, email string) error {
	return DB.Model(&models.User{}).Where("id = ?", userID).Update("email", email).Error
}

func GetUsersByDepartmentAndRole(department, role string) ([]models.User, error) {
	var users []models.User
	err := DB.Where("department = ? AND role = ?", department, role).Find(&users).Error
	return users, err
}
