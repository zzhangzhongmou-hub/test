package cron

import (
	"fmt"
	"test/configs"
	"test/dao"
	"test/service"
	"time"

	"github.com/robfig/cron/v3"
)

var c *cron.Cron

func checkUpcomingDeadlines() {
	now := time.Now()
	deadline := now.Add(24 * time.Hour)

	homeworks, err := dao.GetHomeworksByDeadlineRange(now, deadline)
	if err != nil {
		fmt.Printf("[Cron] 查询作业失败: %v\n", err)
		return
	}
	for _, hw := range homeworks {
		unsubmitted, err := dao.GetUnsubmittedStudents(hw.ID, hw.Department)
		if err != nil {
			fmt.Printf("[Cron] 查询未提交学生失败: %v\n", err)
			continue
		}
		for _, student := range unsubmitted {
			if student.Email == "" {
				continue
			}
			service.EmailSvc.NotifyDeadlineReminder(
				student.Email,
				hw.Title,
				hw.Deadline.Format("2025-01-02 15:04"),
			)
		}
	}
}

func Init() {
	if !configs.Cfg.Cron.Enable {
		fmt.Println("[Cron] 定时任务已禁用")
		return
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	c = cron.New(cron.WithLocation(loc))

	_, err := c.AddFunc(configs.Cfg.Cron.ReminderTime, checkUpcomingDeadlines)
	if err != nil {
		fmt.Printf("[Cron] 添加任务失败: %v\n", err)
		return
	}

	c.Start()
	fmt.Printf("[Cron] 定时任务已启动，执行时间: %s\n", configs.Cfg.Cron.ReminderTime)
}

func Stop() {
	if c != nil {
		c.Stop()
		fmt.Println("[Cron] 定时任务已停止")
	}
}
