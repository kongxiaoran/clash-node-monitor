package mailer

import (
	"fmt"
	"log"
	"mime"
	"net/smtp"
	"strings"
	"time"

	"clash-node-monitor/checker"
	"clash-node-monitor/config"
)

type Mailer struct {
	config config.EmailConfig
	auth   smtp.Auth
}

func NewMailer(config config.EmailConfig) *Mailer {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.SMTPHost)
	return &Mailer{
		config: config,
		auth:   auth,
	}
}

func (m *Mailer) SendAlertEmail(results []checker.CheckResult) error {
	// 筛选出需要告警的节点
	var failedNodes []checker.CheckResult
	for _, result := range results {
		if result.Error != nil && result.ShouldAlert {
			failedNodes = append(failedNodes, result)
		}
	}

	// 如果没有失败的节点，不发送邮件
	if len(failedNodes) == 0 {
		return nil
	}

	// 构建邮件内容
	body := fmt.Sprintf("检测时间: %s\n\n失败节点列表:\n", time.Now().Format("2006-01-02 15:04:05"))
	for _, node := range failedNodes {
		body += fmt.Sprintf("节点: %s\n错误: %v\n\n", node.Name, node.Error)
	}

	// 构建邮件头
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s",
		m.config.From,
		strings.Join(m.config.To, ", "),
		encodeSubject(m.config.Subject),
		body,
	)

	// 发送邮件
	addr := fmt.Sprintf("%s:%d", m.config.SMTPHost, m.config.SMTPPort)
	log.Println("发送告警邮件:\n", body)
	return smtp.SendMail(addr, m.auth, m.config.From, m.config.To, []byte(msg))
}

func encodeSubject(subject string) string {
	// 使用 quoted-printable 编码邮件标题
	encodedSubject := mime.QEncoding.Encode("utf-8", subject)
	return fmt.Sprintf("Subject: =?utf-8?Q?%s?=\n", encodedSubject)
}
