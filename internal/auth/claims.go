package auth

// AppUserClaims are JWT claims issued for a app end-user
// (model.AppUser) authenticating into a App.
type AppUserClaims struct {
	BaseClaims

	AppCode string `json:"app_code"`
}
