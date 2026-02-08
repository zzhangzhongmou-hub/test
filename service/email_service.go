package service

import (
	"fmt"
	"test/dao"
	"test/pkg/email"
)

type EmailService struct{}

// NotifyNewHomework
func (s *EmailService) NotifyNewHomework(department string, title string) {
	users, err := dao.GetUsersByDepartmentAndRole(department, "student")
	if err != nil {
		fmt.Printf("[Email] 查询学生失败: %v\n", err)
		return
	}

	var toList []string
	for _, u := range users {
		if u.Email != "" {
			toList = append(toList, u.Email)
		}
	}

	if len(toList) == 0 {
		fmt.Printf("[Email] 部门 %s 无学生绑定邮箱，跳过通知\n", department)
		return
	}

	subject := fmt.Sprintf("【新作业发布】%s", title)
	body := fmt.Sprintf(`
        <h2>新作业发布通知</h2>
        <p>您所在的部门发布了新作业：<strong>%s</strong></p>
        <p>请及时登录系统查看详情并完成作业。</p>
        <hr>
        <p style="color:#666;font-size:12px;">此邮件由红岩网校作业系统自动发送，请勿回复。</p>
    `, title)

	email.Send(toList, subject, body, true)
}

func (s *EmailService) NotifyReview(studentID uint, homeworkTitle string, score int, comment string) {
	user, err := dao.GetUserByID(studentID)
	if err != nil || user.Email == "" {
		return
	}

	subject := fmt.Sprintf("【作业批改完成】%s", homeworkTitle)
	body := fmt.Sprintf(`
        <h2>作业批改通知</h2>
        <p>您的作业<strong>%s</strong>已被批改</p>
        <p><strong>得分：</strong>%d</p>
        <p><strong>评语：</strong>%s</p>
        <hr>
        <p style="color:#666;font-size:12px;">此邮件由红岩网校作业系统自动发送，请勿回复。</p>
    `, homeworkTitle, score, comment)

	email.Send([]string{user.Email}, subject, body, true)
}

func (s *EmailService) NotifyDeadlineReminder(useremail string, homeworkTitle string, deadline string) {
	subject := fmt.Sprintf("【截止提醒】作业 %s 即将截止", homeworkTitle)
	body := fmt.Sprintf(`
        <h2>作业截止提醒</h2>
        <p>您未完成的作业<strong>%s</strong>将在24小时内截止</p>
        <p><strong>截止时间：</strong>%s</p>
        <p>请尽快登录系统提交作业。</p>
        <hr>
        <p style="color:#666;font-size:12px;">此邮件由红岩网校作业系统自动发送，请勿回复。</p>
    `, homeworkTitle, deadline)

	email.Send([]string{useremail}, subject, body, true)
}

var EmailSvc = &EmailService{}
