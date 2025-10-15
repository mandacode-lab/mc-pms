package configs

type IDPrefixConfig struct {
	Service string `env:"SERVICE" validate:"required"`
	// Client  string `env:"CLIENT" validate:"required"`
}
