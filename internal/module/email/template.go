package email

import (
	"embed"
	"html/template"
)

//go:embed templates/*.html
var templatesFS embed.FS

var templates = template.Must(template.ParseFS(templatesFS, "templates/*.html"))

type TemplateName string

const (
	TemplateActivateAccount  TemplateName = "activate_account"
	TemplateResetPassword    TemplateName = "reset_password"
	TemplateTenantInvitation TemplateName = "tenant_invitation"
)

var subjects = map[TemplateName]string{
	TemplateActivateAccount:  "Activate your account",
	TemplateResetPassword:    "Reset your password",
	TemplateTenantInvitation: "You've been invited",
}

// ActivateAccountData is the data expected by the "activate_account" template.
type ActivateAccountData struct {
	Name string
	Link string
}

// ResetPasswordData is the data expected by the "reset_password" template.
type ResetPasswordData struct {
	Name string
	Link string
}

// TenantInvitationData is the data expected by the "tenant_invitation" template.
type TenantInvitationData struct {
	TenantName  string
	InviterName string
	Link        string
}
