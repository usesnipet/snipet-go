# Email module (SMTP)

Not CRUD, no entity, no HTTP routes — a plain service dependency other modules inject to send transactional email (account activation, password reset, invitation, etc.). Package: `internal/module/email`.

## Config

```go
// config.SMTPConfig (new)
type SMTPConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    From     string // e.g. "Snipet <no-reply@snipet.dev>"
    UseTLS   bool
}
```

## Service

```go
type Service struct {
    cfg config.SMTPConfig
}

func NewService(cfg config.SMTPConfig) *Service

func (s *Service) Send(ctx context.Context, to, subject, htmlBody string) error
func (s *Service) SendTemplate(ctx context.Context, to string, tmpl TemplateName, data any) error
```

`Send` does the raw SMTP dial+auth+deliver (`net/smtp`, or a small lib if TLS/attachments become a pain — implementation detail, not a plan concern). `SendTemplate` renders a Go `html/template` by name and calls `Send`.

## Templates

```go
type TemplateName string

const (
    TemplateActivateAccount   TemplateName = "activate_account"
    TemplateResetPassword     TemplateName = "reset_password"
    TemplateTenantInvitation  TemplateName = "tenant_invitation"
)
```

Each template gets the action link plus context data:

- `activate_account` / `reset_password` — `https://.../activate?token=...` / `.../reset-password?token=...`, plus `Name` (`auth`, see `plan-ee/module/auth.md`).
- `tenant_invitation` — an accept link carrying `tenant_id`, `invitation_id`, and `token` (the frontend page reads these, prompts login/register if needed, then calls `POST /tenants/{tenant_id}/invitations/{id}/accept`), plus `TenantName` and inviter `Name` (`tenant-invitation`, see `plan-ee/module/tenant-invitation.md`).

Template files live under `internal/module/email/templates/*.html`.

## Consumers

`auth` (`plan-ee/module/auth.md`) — `SendTemplate` after issuing an `activate_account`/`reset_password` `Token`. `tenant-invitation` (`plan-ee/module/tenant-invitation.md`) — `SendTemplate` after `Create`.