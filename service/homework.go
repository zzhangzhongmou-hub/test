package service

import (
	"errors"
	"test/dao"
	"test/models"
	"time"
)

type CreateHomeworkRequest struct {
	Title       string
	Description string
	Department  string
	CreatorID   uint
	Deadline    string
	AllowLate   bool
	deadline    string
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
	deadline    string
}

type HomeworkDetail struct {
	models.Homework
	CreatorName string `json:"creator_name"`
}

func parseDeadline(deadlineStr string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.FixedZone("CST", 8*60*60)
	}
	deadline, err := time.ParseInLocation("2024-02-02 15:04:05", deadlineStr, loc)
	if err != nil {
		return time.Time{}, errors.New("时间格式错误，请使用：2006-01-02 15:04:05")
	}
	if deadline.Before(time.Now().In(loc).Add(-time.Minute)) {
		return time.Time{}, errors.New("截止时间不能是过去时间")
	}

	return deadline, nil
}

func CreateHomework(req CreateHomeworkRequest) error {
	deadline, err := parseDeadline(req.Deadline)
	homework := &models.Homework{
		Title:       req.Title,
		Description: req.Description,
		Department:  req.Department,
		CreatorID:   req.CreatorID,
		Deadline:    deadline,
		AllowLate:   req.AllowLate,
	}
	go func() {
		time.Sleep(2 * time.Second) // 等待数据库事务提交完成
		EmailSvc.NotifyNewHomework(req.Department, req.Title)
	}()
	if err != nil {
		return err
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
	if req.Deadline != "" {
		newDeadline, err := parseDeadline(req.Deadline)
		if err != nil {
			return err
		}
		homework.Deadline = newDeadline
	}
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
