package tool

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"gin-skeleton/model/admin/system"

	"github.com/jordan-wright/email"
)

type Email struct {
	Server   string `json:"server"`   // 收信服务器
	Port     string `json:"port"`     // 端口
	Sender   string `json:"sender"`   // 发件人
	Account  string `json:"account"`  // 账号
	Password string `json:"password"` // 密码或授权码
}

// FindEmail 获取邮件配置
func FindEmail() (email *Email, err error) {
	err = system.NewCommonConfig().FindConfigValueTo(system.CommonConfigModuleEmailServer, "", &email)
	return
}

// Send 发送邮件
func (m *Email) Send(subject, body string, toEmails, ccEmails, bccEmails []string) (err error) {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", m.Sender, m.Account)
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
