package configs

type SystemClientConfig struct {
	SystemServiceName  string `env:"SYSTEM_SERVICE_NAME" envDefault:"system" validate:"required"`
	SystemClientName   string `env:"SYSTEM_CLIENT_NAME" envDefault:"system-client" validate:"required"`
	SystemClientID     string `env:"SYSTEM_CLIENT_ID" validate:"required,uuid4"`
	SystemClientSecret string `env:"SYSTEM_CLIENT_SECRET" validate:"required,min=32"`
}
