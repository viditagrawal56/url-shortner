package email

import (
	"fmt"
	"log"
	"net/smtp"
)

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
}

type Service struct {
	cfg *EmailConfig
}

func NewEmailService(cfg *EmailConfig) *Service {
	return &Service{
		cfg: cfg,
	}
}

func (s *Service) SendMagicLink(to, shortCode, token, baseURL string) error {
	auth := smtp.PlainAuth("", s.cfg.SMTPUsername, s.cfg.SMTPPassword, s.cfg.SMTPHost)

	log.Printf("Attempting to send email via SMTP:\nHost: %s\nPort: %d\nUsername: %s\nFrom: %s\nTo: %s",
		s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPUsername, s.cfg.FromEmail, to)

	verificationLink := fmt.Sprintf("%s/verify?token=%s", baseURL, token)

	subject := "Your Magic Link for URL Access"
	body := fmt.Sprintf(`
	<html>
	<body>
		<h2>Access Your Requested URL</h2>
		<p>You requested access to a shortened URL. Click the link below to proceed:</p>
		<p><a href="%s">Click here to access your requested URL</a></p>
		<p>This link will expire in 15 minutes.</p>
		<p>If you didn't request this, you can safely ignore this email.</p>
	</body>
	</html>
	`, verificationLink)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := fmt.Sprintf("Subject: %s\n%s\n%s", subject, mime, body)

	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)

	if err := smtp.SendMail(addr, auth, s.cfg.FromEmail, []string{to}, []byte(msg)); err != nil {
		log.Printf("Failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *Service) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.cfg.SMTPUsername, s.cfg.SMTPPassword, s.cfg.SMTPHost)

	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	msg := fmt.Sprintf("Subject: %s\n%s\n%s", subject, mime, body)

	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)

	return smtp.SendMail(addr, auth, s.cfg.FromEmail, []string{to}, []byte(msg))
}
