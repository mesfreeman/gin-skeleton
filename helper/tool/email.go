package tool

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type Email struct {
	Server   string   `json:"server"`   // 收信服务器
	Port     string   `json:"port"`     // 端口
	Fromer   string   `json:"fromer"`   // 发件人
	Account  string   `json:"account"`  // 账号
	Password string   `json:"password"` // 密码或授权码
	ToEmails []string `json:"toEmails"` // 收件人邮箱
}

// Send 发送邮件
func (m *Email) Send(subject, body string, toEmails, ccEmails, bccEmails []string) (err error) {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", m.Fromer, m.Account)
	e.To = toEmails
	e.Cc = ccEmails
	e.Bcc = bccEmails
	e.Subject = subject
	e.HTML = []byte(body)

	address := fmt.Sprintf("%s:%s", m.Server, m.Port)
	auth := smtp.PlainAuth("", m.Account, m.Password, m.Server)
	err = e.SendWithTLS(address, auth, &tls.Config{ServerName: m.Server})
	return
}

// ConfigNoticer 邮件配置通知
func (m *Email) ConfigNoticer(opertor string) (err error) {
	subject := "邮件配置变更通知"
	body := fmt.Sprintf("邮件配置变更成功，操作人：%s", opertor)
	err = m.Send(subject, body, m.ToEmails, []string{}, []string{})
	return
}
