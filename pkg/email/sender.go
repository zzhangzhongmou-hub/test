package email

import (
	"fmt"
	"test/configs"

	"github.com/go-gomail/gomail"
)

type MailMessage struct {
	To      []string
	Subject string
	Body    string
	IsHTML  bool
}

var (
	mailQueue chan *MailMessage
	dialer    *gomail.Dialer
)

func Init() {
	cfg := configs.Cfg.SMTP
	if !cfg.Enable {
		fmt.Println("[Email] 邮件功能已禁用")
		return
	}

	dialer = gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)

	if cfg.Port == 465 {
		dialer.SSL = true
	}

	mailQueue = make(chan *MailMessage, 100)

	go sendWorker()

	fmt.Printf("[Email] 邮件服务已启动 (%s:%d)\n", cfg.Host, cfg.Port)
}

// email.Send([]string{"user@qq.com"}, "标题", "内容", true)
func Send(to []string, subject, body string, isHTML bool) {
	if !configs.Cfg.SMTP.Enable || mailQueue == nil {
		fmt.Printf("[Email] 跳过发送（功能未启用）: %v\n", to)
		return
	}

	msg := &MailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
		IsHTML:  isHTML,
	}

	select {
	case mailQueue <- msg:
		fmt.Printf("[Email] 已入队，收件人: %v\n", to)
	default:
		fmt.Printf("[Email] 队列已满，丢弃邮件: %v\n", to)
	}
}

func sendWorker() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[Email] Worker崩溃: %v，正在重启...\n", r)
			go sendWorker()
		}
	}()

	for msg := range mailQueue {
		if err := sendEmail(msg); err != nil {
			fmt.Printf("[Email] 发送失败 %v: %v\n", msg.To, err)
		} else {
			fmt.Printf("[Email] 发送成功 %v\n", msg.To)
		}
	}
}

func sendEmail(msg *MailMessage) error {
	cfg := configs.Cfg.SMTP

	m := gomail.NewMessage()
	m.SetHeader("From", cfg.Username)
	m.SetHeader("To", msg.To...)
	m.SetHeader("Subject", msg.Subject)

	if msg.IsHTML {
		m.SetBody("text/html", msg.Body)
	} else {
		m.SetBody("text/plain", msg.Body)
	}

	return dialer.DialAndSend(m)
}

func Close() {
	if mailQueue != nil {
		close(mailQueue)
		fmt.Println("[Email] 邮件服务已关闭")
	}
}
