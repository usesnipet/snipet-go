package auth

// ClientUserClaims are JWT claims issued for a client-widget end-user
// (model.ClientUser) authenticating into a Client.
type ClientUserClaims struct {
	BaseClaims

	ClientCode string `json:"client_code"`
}

// PlatformUserClaims are JWT claims issued for a tenant-staff user.
type PlatformUserClaims struct {
	BaseClaims
}
