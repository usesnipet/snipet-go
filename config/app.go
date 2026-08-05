package config

type AppConfig struct {
	InheritClient     bool   `env:"INHERIT_CLIENT, default=true"`
	InheritClientName string `env:"INHERIT_CLIENT_NAME, default=Snipet"`
	InheritClientCode string `env:"INHERIT_CLIENT_CODE, default=default"`
}
