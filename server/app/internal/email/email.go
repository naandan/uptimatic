package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/smtp"
	"time"
	"uptimatic/internal/config"

	"github.com/jordan-wright/email"
)

type EmailTask struct {
	cfg      *config.Config
	tplCache map[EmailType]*template.Template
}

type EmailType string

const (
	EmailWelcome       EmailType = "welcome"
	EmailVerify        EmailType = "verify"
	EmailPasswordReset EmailType = "password_reset"
	EmailDown          EmailType = "down"
)

type EmailPayload struct {
	To      string         `json:"to"`
	Subject string         `json:"subject"`
	Type    EmailType      `json:"type"`
	Data    map[string]any `json:"data"`
}

func NewEmailTask(cfg *config.Config) (*EmailTask, error) {
	t := &EmailTask{cfg: cfg, tplCache: map[EmailType]*template.Template{}}

	types := []EmailType{EmailWelcome, EmailVerify, EmailPasswordReset, EmailDown}
	for _, typ := range types {
		tpl, err := template.ParseFiles(fmt.Sprintf("internal/email/templates/%s.html", typ))
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", typ, err)
		}
		t.tplCache[typ] = tpl
	}
	return t, nil
}

func (m *EmailTask) SendEmail(ctx context.Context, to, subject string, emailType EmailType, data map[string]any) error {
	tpl, ok := m.tplCache[emailType]
	if !ok {
		return fmt.Errorf("template not found for type %s", emailType)
	}

	var body bytes.Buffer
	if err := tpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	e := email.NewEmail()
	e.From = m.cfg.EmailFrom
	e.To = []string{to}
	e.Subject = subject
	e.Text = []byte("Please view this email in HTML format")
	e.HTML = body.Bytes()

	addr := fmt.Sprintf("%s:%d", m.cfg.EmailSmtpHost, m.cfg.EmailSmtpPort)
	auth := smtp.PlainAuth("", m.cfg.EmailSmtpUser, m.cfg.EmailSmtpPass, m.cfg.EmailSmtpHost)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- e.Send(addr, auth) }()

	select {
	case <-ctx.Done():
		return fmt.Errorf("email send timed out")
	case err := <-done:
		if err != nil {
			return err
		}
	}

	return nil
}
