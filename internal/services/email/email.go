package email

import (
	"fmt"
	"github.com/Gitrupesh20/real-time-notification-system/cmd/config"
	"net/smtp"
	"strings"
)

type EmailService struct {
	config   config.Config
	username string
	password string
	smtpHost string
	smtpPort string
	auth     smtp.Auth
}

type EmailFormat struct {
	To      string
	CC      []string
	Subject string
	Body    string
}

func NewEmailService(config config.Config) *EmailService {
	es := &EmailService{
		config:   config,
		username: "rursharma02@gmail.com",
		password: "rursharma02",
		smtpHost: "smtp.gmail.com",
		smtpPort: "597",
	}

	return es
}

func (e *EmailService) StartEmailServer() {
	//authontication
	auth := smtp.PlainAuth("", e.username, e.password, e.smtpHost)
	e.auth = auth
}

func (e *EmailService) SendEmail(format EmailFormat) {
	recipients := []string{format.To}
	if len(format.CC) > 0 {
		recipients = append(recipients, format.CC...)
	}

	body := []byte(fmt.Sprintf("FROM: %s\r\nTo: %s\r\n", e.username, format.To))
	body = append(body, []byte(fmt.Sprintf("Cc: %s\r\n", strings.Join(format.CC, " ")))...)

}
