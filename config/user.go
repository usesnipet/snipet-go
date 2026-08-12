package config

type UserConfig struct {
	AdminName     string `env:"ADMIN_NAME, default=Admin"`
	AdminEmail    string `env:"ADMIN_EMAIL, default=admin@admin.com"`
	AdminPassword string `env:"ADMIN_PASSWORD"`
}
