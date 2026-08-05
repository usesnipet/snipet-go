package manager

import "github.com/usesnipet/snipet/pkg/jsonx"

type Configuration struct {
	Key    string        `json:"key"`
	Config jsonx.JSONMap `json:"config"`
}
