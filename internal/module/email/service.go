package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"

	"github.com/usesnipet/snipet/config"
	"github.com/usesnipet/snipet/internal/logger"
)

type Service struct {
	cfg    config.SMTPConfig
	logger *logger.Logger
}

func NewService(cfg config.SMTPConfig, logger *logger.Logger) *Service {
	return &Service{cfg: cfg, logger: logger}
}

// SendTemplate renders the named template with data and sends it to to.
func (s *Service) SendTemplate(ctx context.Context, to string, tmpl TemplateName, data any) error {
	if !s.cfg.Enable {
		s.logger.Infof("email: not enabled, skipping send, to: %s, template: %s, data: %v", to, tmpl, data)
		return nil
	}

	var body bytes.Buffer
	if err := templates.ExecuteTemplate(&body, string(tmpl)+".html", data); err != nil {
		s.logger.Errorf("email: render template %q: %v", tmpl, err)
		return ErrEmailSendingFailed
	}

	return s.Send(ctx, to, subjects[tmpl], body.String())
}

// Send dials the configured SMTP server and delivers a single HTML email.
func (s *Service) Send(ctx context.Context, to, subject, htmlBody string) error {
	if !s.cfg.Enable {
		s.logger.Infof("email: not enabled, skipping send, to: %s, subject: %s", to, subject)
		return nil
	}

	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		s.logger.Errorf("email: dial %s: %v", addr, err)
		return ErrEmailSendingFailed
	}

	if s.cfg.UseTLS {
		conn = tls.Client(conn, &tls.Config{ServerName: s.cfg.Host})
	}

	client, err := smtp.NewClient(conn, s.cfg.Host)
	if err != nil {
		_ = conn.Close()
		s.logger.Errorf("email: new smtp client: %v", err)
		return ErrEmailSendingFailed
	}
	defer client.Close()

	if !s.cfg.UseTLS {
		if ok, _ := client.Extension("STARTTLS"); ok {
			if err := client.StartTLS(&tls.Config{ServerName: s.cfg.Host}); err != nil {
				s.logger.Errorf("email: starttls: %v", err)
				return ErrEmailSendingFailed
			}
		}
	}

	if s.cfg.Username != "" {
		auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
		if err := client.Auth(auth); err != nil {
			s.logger.Errorf("email: auth: %v", err)
			return ErrEmailSendingFailed
		}
	}

	if err := client.Mail(s.cfg.From); err != nil {
		s.logger.Errorf("email: mail from: %v", err)
		return ErrEmailSendingFailed
	}
	if err := client.Rcpt(to); err != nil {
		s.logger.Errorf("email: rcpt to: %v", err)
		return ErrEmailSendingFailed
	}

	w, err := client.Data()
	if err != nil {
		s.logger.Errorf("email: data: %v", err)
		return ErrEmailSendingFailed
	}

	if _, err := w.Write(buildMessage(s.cfg.From, to, subject, htmlBody)); err != nil {
		_ = w.Close()
		s.logger.Errorf("email: write body: %v", err)
		return ErrEmailSendingFailed
	}
	if err := w.Close(); err != nil {
		s.logger.Errorf("email: close body: %v", err)
		return ErrEmailSendingFailed
	}

	return client.Quit()
}

func buildMessage(from, to, subject, htmlBody string) []byte {
	var msg bytes.Buffer
	fmt.Fprintf(&msg, "From: %s\r\n", from)
	fmt.Fprintf(&msg, "To: %s\r\n", to)
	fmt.Fprintf(&msg, "Subject: %s\r\n", subject)
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)
	return msg.Bytes()
}
