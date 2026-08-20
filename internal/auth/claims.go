package auth

// AppUserClaims are JWT claims issued for a app end-user
// (model.AppUser) authenticating into a App.
type AppUserClaims struct {
	BaseClaims

	AppCode string `json:"app_code"`
}

// PlatformUserClaims are JWT claims issued for a tenant-staff user.
type PlatformUserClaims struct {
	BaseClaims
}
