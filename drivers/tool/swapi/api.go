package swapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/usesnipet/snipet/pkg/driver/tool"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

const baseURL = "https://swapi.info/api"

var httpClient = &http.Client{Timeout: 15 * time.Second}

// NewAPI returns a tool.API that talks to https://swapi.info/api.
func NewAPI() tool.API {
	return tool.API{
		TestConnection: testConnection,
		Call:           call,
	}
}

func testConnection(ctx context.Context, _ jsonx.JSONMap) error {
	_, err := get(ctx, baseURL)
	if err != nil {
		return fmt.Errorf("swapi connection failed: %w", err)
	}
	return nil
}

func call(ctx context.Context, tc tool.Call) (tool.Result, error) {
	path, err := resolvePath(tc)
	if err != nil {
		return tool.Result{}, err
	}

	body, err := get(ctx, baseURL+path)
	if err != nil {
		return tool.Result{}, err
	}

	return tool.Result{
		Tool:      tc.Tool,
		Arguments: tc.Arguments,
		Result:    body,
	}, nil
}

func resolvePath(tc tool.Call) (string, error) {
	switch tc.Tool {
	case "list_people":
		return "/people", nil
	case "get_people":
		return resourcePath("people", tc.Arguments)
	case "list_films":
		return "/films", nil
	case "get_films":
		return resourcePath("films", tc.Arguments)
	case "list_planets":
		return "/planets", nil
	case "get_planets":
		return resourcePath("planets", tc.Arguments)
	case "list_species":
		return "/species", nil
	case "get_species":
		return resourcePath("species", tc.Arguments)
	case "list_starships":
		return "/starships", nil
	case "get_starships":
		return resourcePath("starships", tc.Arguments)
	case "list_vehicles":
		return "/vehicles", nil
	case "get_vehicles":
		return resourcePath("vehicles", tc.Arguments)
	default:
		return "", fmt.Errorf("unknown tool %q", tc.Tool)
	}
}

func resourcePath(resource string, args map[string]any) (string, error) {
	id, err := argumentID(args)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/%s/%d", resource, id), nil
}

func argumentID(args map[string]any) (int, error) {
	raw, ok := args["id"]
	if !ok {
		return 0, fmt.Errorf("id is required")
	}

	switch v := raw.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		id, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("invalid id %q: %w", v, err)
		}
		return id, nil
	default:
		return 0, fmt.Errorf("invalid id type %T", raw)
	}
}

func get(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "snipet-swapi-driver/1.0")

	res, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request %s: %w", url, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("swapi %s returned %d: %s", url, res.StatusCode, string(body))
	}

	return string(body), nil
}
