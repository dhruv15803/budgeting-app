package email

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/dhruv15803/budgeting-app/internal/config"
)

type Sender struct {
	cfg config.SMTPConfig
}

func NewSender(cfg config.SMTPConfig) *Sender {
	return &Sender{cfg: cfg}
}

func (s *Sender) SendVerification(to string, verificationURL string) error {
	if s.cfg.Host == "" {
		return fmt.Errorf("SMTP is not configured (SMTP_HOST empty)")
	}
	from := s.cfg.From
	if from == "" {
		return fmt.Errorf("SMTP_FROM is not set")
	}
	subject := "Verify your email address"
	body := strings.Builder{}
	body.WriteString("Verify your account by opening this link:\n\n")
	body.WriteString(verificationURL)
	body.WriteString("\n")
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	auth := smtp.PlainAuth("", s.cfg.User, s.cfg.Password, s.cfg.Host)
	msg := []byte("To: " + to + "\r\nSubject: " + subject + "\r\n\r\n" + body.String() + "\r\n")
	return smtp.SendMail(addr, auth, from, []string{to}, msg)
}
