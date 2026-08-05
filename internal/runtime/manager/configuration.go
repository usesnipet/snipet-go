package manager

import "github.com/usesnipet/snipet/internal/util"

type Configuration struct {
	Key    string       `json:"key"`
	Config util.JSONMap `json:"config"`
}
