package config

type AppConfig struct {
	InheritClient bool `env:"INHERIT_CLIENT, default=true"`
}
