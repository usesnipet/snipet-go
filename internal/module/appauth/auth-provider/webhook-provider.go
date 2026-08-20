package auth_provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	apperr "github.com/usesnipet/snipet/internal/app-err"
	"github.com/usesnipet/snipet/internal/model"
)

const WebhookProviderName ProviderName = "webhook"

var excludedWebhookHeaders = map[string]struct{}{
	"Host":                {},
	"Connection":          {},
	"Keep-Alive":          {},
	"Proxy-Authenticate":  {},
	"Proxy-Authorization": {},
	"Te":                  {},
	"Trailer":             {},
	"Transfer-Encoding":   {},
	"Upgrade":             {},
}

type WebhookProvider struct {
	client *http.Client
}

func NewWebhookProvider() IProvider {
	return &WebhookProvider{
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (p *WebhookProvider) Name() ProviderName {
	return WebhookProviderName
}

func (p *WebhookProvider) Validate(
	ctx context.Context,
	appCode string,
	authConfig *model.AppAuthConfig,
) error {
	if authConfig.Webhook.URL == "" {
		return apperr.BadRequest("app webhook url not configured")
	}
	if !authConfig.Webhook.Enabled {
		return apperr.BadRequest("app webhook is not enabled")
	}
	return nil
}

func (p *WebhookProvider) Authenticate(
	ctx context.Context,
	appCode string,
	authConfig *model.AppAuthConfig,
	req *http.Request,
) (*Identity, error) {
	webhookURL, err := url.Parse(authConfig.Webhook.URL)
	if err != nil {
		return nil, apperr.BadRequest("invalid app webhook url")
	}

	query := webhookURL.Query()
	for key, values := range req.URL.Query() {
		for _, value := range values {
			query.Add(key, value)
		}
	}
	webhookURL.RawQuery = query.Encode()

	webhookReq, err := http.NewRequestWithContext(ctx, req.Method, webhookURL.String(), req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create webhook request: %w", err)
	}

	for key, values := range req.Header {
		if _, excluded := excludedWebhookHeaders[http.CanonicalHeaderKey(key)]; excluded {
			continue
		}
		for _, value := range values {
			webhookReq.Header.Add(key, value)
		}
	}

	webhookReq.Header.Set("User-Agent", "Snipet-Webhook/1.0")

	resp, err := p.client.Do(webhookReq)
	if err != nil {
		return nil, apperr.NetworkError("failed to call app webhook")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, apperr.Unauthorized("app webhook rejected authentication")
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, apperr.NetworkError("app webhook returned an error")
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		return nil, apperr.BadRequest("app webhook response is not a JSON")
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, apperr.NetworkError("failed to read app webhook response")
	}

	var identity Identity
	if err := json.Unmarshal(body, &identity); err != nil {
		return nil, apperr.BadRequest("invalid app webhook response")
	}
	return &identity, nil
}
