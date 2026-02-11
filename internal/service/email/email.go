package email

import (
	"fmt"
	"net/smtp"
	"strings"

	"week5/internal/config"
)

// Notifier sends emails asynchronously using smtp.SendMail.
type Notifier struct {
	host string
	port int
	user string
	pass string
	to   string
}

func NewNotifier(cfg *config.Config) *Notifier {
	return &Notifier{
		host: cfg.SMTPHost,
		port: cfg.SMTPPort,
		user: cfg.SMTPUser,
		pass: cfg.SMTPPass,
		to:   cfg.SMTPTo,
	}
}

func (n *Notifier) SendAsync(subject, body string) {
	go func() {
		auth := smtp.PlainAuth("", n.user, n.pass, n.host)
		addr := fmt.Sprintf("%s:%d", n.host, n.port)
		msg := buildMessage(n.user, n.to, subject, body)
		_ = smtp.SendMail(addr, auth, n.user, []string{n.to}, []byte(msg))
	}()
}

func buildMessage(from, to, subject, body string) string {
	headers := []string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=\"utf-8\"",
		"",
		body,
	}
	return strings.Join(headers, "\r\n")
}
