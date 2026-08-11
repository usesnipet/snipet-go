package clientauth

import (
	coreauth "github.com/usesnipet/snipet/internal/auth"
)

// UserClaims are the JWT claims issued for a client-widget end-user
// (model.ClientUser) authenticating into a Client.
type UserClaims struct {
	coreauth.BaseClaims

	ClientCode string `json:"client_code"`
}
