package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"

	"github.com/usesnipet/snipet/config"
)

type Service struct {
	cfg config.SMTPConfig
}

func NewService(cfg config.SMTPConfig) *Service {
	return &Service{cfg: cfg}
}

// SendTemplate renders the named template with data and sends it to to.
func (s *Service) SendTemplate(ctx context.Context, to string, tmpl TemplateName, data any) error {
	var body bytes.Buffer
	if err := templates.ExecuteTemplate(&body, string(tmpl)+".html", data); err != nil {
		return fmt.Errorf("email: render template %q: %w", tmpl, err)
	}

	return s.Send(ctx, to, subjects[tmpl], body.String())
}

// Send dials the configured SMTP server and delivers a single HTML email.
func (s *Service) Send(ctx context.Context, to, subject, htmlBody string) error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("email: dial %s: %w", addr, err)
	}

	if s.cfg.UseTLS {
		conn = tls.Client(conn, &tls.Config{ServerName: s.cfg.Host})
	}

	client, err := smtp.NewClient(conn, s.cfg.Host)
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("email: new smtp client: %w", err)
	}
	defer client.Close()

	if !s.cfg.UseTLS {
		if ok, _ := client.Extension("STARTTLS"); ok {
			if err := client.StartTLS(&tls.Config{ServerName: s.cfg.Host}); err != nil {
				return fmt.Errorf("email: starttls: %w", err)
			}
		}
	}

	if s.cfg.Username != "" {
		auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("email: auth: %w", err)
		}
	}

	if err := client.Mail(s.cfg.From); err != nil {
		return fmt.Errorf("email: mail from: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("email: rcpt to: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("email: data: %w", err)
	}

	if _, err := w.Write(buildMessage(s.cfg.From, to, subject, htmlBody)); err != nil {
		_ = w.Close()
		return fmt.Errorf("email: write body: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("email: close body: %w", err)
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
