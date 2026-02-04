package service

import (
	"errors"
	"test/dao"
	"test/models"
)

type CreateHomeworkRequest struct {
	Title       string
	Description string
	Department  string
	CreatorID   uint
	Deadline    string
	AllowLate   bool
}

type UpdateHomeworkRequest struct {
	ID          uint
	Title       string
	Description string
	Deadline    string
	AllowLate   bool
	Version     int
	UpdaterID   uint
	UpdaterDept string
}

type HomeworkDetail struct {
	models.Homework
	CreatorName string `json:"creator_name"`
}

func CreateHomework(req CreateHomeworkRequest) error {
	homework := &models.Homework{
		Title:       req.Title,
		Description: req.Description,
		Department:  req.Department,
		CreatorID:   req.CreatorID,
		AllowLate:   req.AllowLate,
		// Deadline 解析暂时用当前时间，明天完善时间解析
	}
	return dao.CreateHomework(homework)
}

func GetHomeworkList(department string, page, pageSize int) ([]models.Homework, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return dao.GetHomeworksByDepartment(department, page, pageSize)
}

func GetHomeworkDetail(id uint) (*models.Homework, error) {
	return dao.GetHomeworkByID(id)
}

func UpdateHomework(req UpdateHomeworkRequest) error {
	homework, err := dao.GetHomeworkByID(req.ID)
	if err != nil {
		return errors.New("作业不存在")
	}
	if homework.Department != req.UpdaterDept {
		return errors.New("只能修改本部门的作业")
	}
	if homework.Version != req.Version {
		return errors.New("作业已被其他人修改，请刷新后重试")
	}
	homework.Title = req.Title
	homework.Description = req.Description
	// homework.Deadline = ... // 明天完善时间解析
	homework.AllowLate = req.AllowLate
	homework.Version = req.Version

	return dao.UpdateHomework(homework)
}

func DeleteHomework(id uint, operatorDept string) error {
	homework, err := dao.GetHomeworkByID(id)
	if err != nil {
		return errors.New("作业不存在")
	}
	if homework.Department != operatorDept {
		return errors.New("只能删除本部门的作业")
	}

	return dao.DeleteHomework(id)
}
