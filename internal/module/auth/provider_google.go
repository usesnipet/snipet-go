package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/usesnipet/snipet/config"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"golang.org/x/oauth2"
	googleoauth "golang.org/x/oauth2/google"
)

type GoogleProvider struct {
	oauthConfig *oauth2.Config
}

func NewGoogleProvider(cfg config.AuthConfig) IProvider {
	return &GoogleProvider{
		oauthConfig: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     googleoauth.Endpoint,
		},
	}
}

func (p *GoogleProvider) Name() ProviderName {
	return ProviderGoogle
}

func (p *GoogleProvider) AuthorizationURL(state string) string {
	return p.oauthConfig.AuthCodeURL(state)
}

type googleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func (p *GoogleProvider) Exchange(ctx context.Context, code string) (*Identity, error) {
	token, err := p.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, apperr.Unauthorized("invalid google authorization code")
	}

	client := p.oauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, apperr.BadRequest("failed to fetch google user info")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apperr.BadRequest("failed to fetch google user info")
	}

	var info googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, apperr.BadRequest("invalid google user info response")
	}

	return &Identity{
		ExternalID: info.ID,
		Email:      info.Email,
		Name:       info.Name,
		Picture:    info.Picture,
	}, nil
}
