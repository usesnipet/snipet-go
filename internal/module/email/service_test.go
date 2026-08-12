package email_test

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/usesnipet/snipet/config"
	"github.com/usesnipet/snipet/internal/logger"
	"github.com/usesnipet/snipet/internal/module/email"
)

// fakeSMTPServer is a minimal SMTP server for tests — no auth/STARTTLS
// support, just enough of the protocol (EHLO/MAIL/RCPT/DATA/QUIT) to accept
// a single message per connection and record its raw bytes.
type fakeSMTPServer struct {
	listener net.Listener

	mu       sync.Mutex
	messages []string
}

func startFakeSMTPServer(t *testing.T) *fakeSMTPServer {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	srv := &fakeSMTPServer{listener: ln}
	go srv.serve()
	t.Cleanup(func() { _ = ln.Close() })

	return srv
}

func (s *fakeSMTPServer) config() config.SMTPConfig {
	host, portStr, _ := net.SplitHostPort(s.listener.Addr().String())
	port, _ := strconv.Atoi(portStr)
	return config.SMTPConfig{
		Host: host,
		Port: port,
		From: "Snipet <no-reply@snipet.dev>",
	}
}

func (s *fakeSMTPServer) lastMessage() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.messages) == 0 {
		return ""
	}
	return s.messages[len(s.messages)-1]
}

func (s *fakeSMTPServer) serve() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		go s.handle(conn)
	}
}

func (s *fakeSMTPServer) handle(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writeLine := func(line string) { fmt.Fprintf(conn, "%s\r\n", line) }
	writeLine("220 localhost SMTP fake")

	var inData bool
	var data strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimRight(line, "\r\n")

		if inData {
			if line == "." {
				inData = false
				s.mu.Lock()
				s.messages = append(s.messages, data.String())
				s.mu.Unlock()
				data.Reset()
				writeLine("250 OK")
				continue
			}
			data.WriteString(line)
			data.WriteString("\n")
			continue
		}

		switch upper := strings.ToUpper(line); {
		case strings.HasPrefix(upper, "EHLO"), strings.HasPrefix(upper, "HELO"):
			writeLine("250 localhost")
		case strings.HasPrefix(upper, "MAIL FROM"), strings.HasPrefix(upper, "RCPT TO"):
			writeLine("250 OK")
		case upper == "DATA":
			inData = true
			writeLine("354 End data with <CR><LF>.<CR><LF>")
		case upper == "QUIT":
			writeLine("221 Bye")
			return
		default:
			writeLine("250 OK")
		}
	}
}

func TestSendDeliversHTMLMessage(t *testing.T) {
	t.Parallel()

	srv := startFakeSMTPServer(t)
	svc := email.NewService(srv.config(), logger.NewLogger(logger.LevelDebug))

	err := svc.Send(context.Background(), "user@example.com", "Hello", "<p>Hi</p>")
	require.NoError(t, err)

	msg := srv.lastMessage()
	assert.Contains(t, msg, "Subject: Hello")
	assert.Contains(t, msg, "To: user@example.com")
	assert.Contains(t, msg, "<p>Hi</p>")
}

func TestSendReturnsErrorWhenServerUnreachable(t *testing.T) {
	t.Parallel()

	svc := email.NewService(config.SMTPConfig{Host: "127.0.0.1", Port: 1, From: "no-reply@snipet.dev"}, logger.NewLogger(logger.LevelDebug))

	err := svc.Send(context.Background(), "user@example.com", "Hello", "<p>Hi</p>")
	require.Error(t, err)
}

func TestSendTemplateRendersAndDeliversEachTemplate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		tmpl            email.TemplateName
		data            any
		expectedSubject string
		expectedBody    []string
	}{
		{
			name:            "activate account",
			tmpl:            email.TemplateActivateAccount,
			data:            email.ActivateAccountData{Name: "Ana", Link: "https://snipet.dev/activate?token=abc"},
			expectedSubject: "Subject: Activate your account",
			expectedBody:    []string{"Hi Ana,", "https://snipet.dev/activate?token=abc"},
		},
		{
			name:            "reset password",
			tmpl:            email.TemplateResetPassword,
			data:            email.ResetPasswordData{Name: "Ana", Link: "https://snipet.dev/reset?token=abc"},
			expectedSubject: "Subject: Reset your password",
			expectedBody:    []string{"Hi Ana,", "https://snipet.dev/reset?token=abc"},
		},
		{
			name: "tenant invitation",
			tmpl: email.TemplateTenantInvitation,
			data: email.TenantInvitationData{
				TenantName:  "Acme",
				InviterName: "Bob",
				Link:        "https://snipet.dev/invite?token=abc",
			},
			expectedSubject: "Subject: You've been invited",
			expectedBody:    []string{"Acme", "Bob invited you", "https://snipet.dev/invite?token=abc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := startFakeSMTPServer(t)
			svc := email.NewService(srv.config(), logger.NewLogger(logger.LevelDebug))

			err := svc.SendTemplate(context.Background(), "user@example.com", tt.tmpl, tt.data)
			require.NoError(t, err)

			msg := srv.lastMessage()
			assert.Contains(t, msg, tt.expectedSubject)
			for _, want := range tt.expectedBody {
				assert.Contains(t, msg, want)
			}
		})
	}
}
