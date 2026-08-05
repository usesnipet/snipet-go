package store

import "fmt"

type Config struct {
	URL    string `json:"url"`
	Table  string `json:"table"`
	Length int    `json:"length"`
}

func (c Config) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("store: url is required")
	}
	if c.Table == "" {
		return fmt.Errorf("store: table is required")
	}
	if err := validateIdent(c.Table); err != nil {
		return fmt.Errorf("store: table: %w", err)
	}
	if c.Length <= 0 {
		return fmt.Errorf("store: length must be greater than zero")
	}
	return nil
}

func validateIdent(name string) error {
	if name == "" {
		return fmt.Errorf("empty identifier")
	}
	for i, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r == '_':
			continue
		case r >= '0' && r <= '9':
			if i == 0 {
				return fmt.Errorf("invalid identifier %q", name)
			}
			continue
		default:
			return fmt.Errorf("invalid identifier %q", name)
		}
	}
	return nil
}
