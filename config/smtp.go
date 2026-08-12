package config

type SMTPConfig struct {
	Host     string `env:"HOST, default=localhost"`
	Port     int    `env:"PORT, default=587"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	From     string `env:"FROM, default=Snipet <no-reply@snipet.dev>"`
	UseTLS   bool   `env:"USE_TLS, default=false"`
	Enable   bool   `env:"ENABLE, default=false"`
}
