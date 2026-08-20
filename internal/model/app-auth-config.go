package model

// AppAuthConfig holds how an App federates its end-users' identity into
// Snipet (see internal/module/clientauth/auth-provider) — never credentials,
// just enough to trust an external_id asserted by the App.
type AppAuthConfig struct {
	OIDC struct {
		Issuer   string `json:"issuer" validate:"omitempty,url"`
		Audience string `json:"audience" validate:"omitempty"`
		Enabled  bool   `json:"enabled"`
	} `json:"oidc"`
	Webhook struct {
		URL     string `json:"url" validate:"omitempty,url"`
		Enabled bool   `json:"enabled"`
	} `json:"webhook"`
}
