package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/usesnipet/snipet/config"
	apperr "github.com/usesnipet/snipet/internal/app-err"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

type GithubProvider struct {
	oauthConfig *oauth2.Config
}

func NewGithubProvider(cfg config.AuthConfig) IProvider {
	return &GithubProvider{
		oauthConfig: &oauth2.Config{
			ClientID:     cfg.GithubClientID,
			ClientSecret: cfg.GithubClientSecret,
			RedirectURL:  cfg.GithubRedirectURL,
			Scopes:       []string{"read:user", "user:email"},
			Endpoint:     githuboauth.Endpoint,
		},
	}
}

func (p *GithubProvider) Name() ProviderName {
	return ProviderGithub
}

func (p *GithubProvider) AuthorizationURL(state string) string {
	return p.oauthConfig.AuthCodeURL(state)
}

type githubUser struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type githubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func (p *GithubProvider) Exchange(ctx context.Context, code string) (*Identity, error) {
	token, err := p.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, apperr.Unauthorized("invalid github authorization code")
	}

	client := p.oauthConfig.Client(ctx, token)

	user, err := fetchGithubUser(client)
	if err != nil {
		return nil, err
	}

	email := user.Email
	if email == "" {
		email, err = fetchGithubPrimaryEmail(client)
		if err != nil {
			return nil, err
		}
	}

	return &Identity{
		ExternalID: strconv.FormatInt(user.ID, 10),
		Email:      email,
		Name:       user.Name,
		Picture:    user.AvatarURL,
	}, nil
}

func fetchGithubUser(client *http.Client) (*githubUser, error) {
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, apperr.BadRequest("failed to fetch github user")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apperr.BadRequest("failed to fetch github user")
	}

	var user githubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, apperr.BadRequest("invalid github user response")
	}
	return &user, nil
}

func fetchGithubPrimaryEmail(client *http.Client) (string, error) {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return "", apperr.BadRequest("failed to fetch github user emails")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", apperr.BadRequest("failed to fetch github user emails")
	}

	var emails []githubEmail
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", apperr.BadRequest("invalid github user emails response")
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}
	return "", apperr.BadRequest("no verified primary github email")
}
