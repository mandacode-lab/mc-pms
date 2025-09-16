package configs

type ManagementConfig struct {
	Env      string         `env:"ENV" envDefault:"prod" validate:"oneof=dev staging prod"`
	KEK      string         `env:"KEK" validate:"required,len=64,hexadecimal"` // 256-bit hex-encoded key (64 chars)
	Server   ServerConfig   `envPrefix:""`
	Postgres PostgresConfig `envPrefix:""`
	Redis    RedisConfig    `envPrefix:""`
	IAM      IAMConfig      `envPrefix:"IAM_"`
}

type IAMConfig struct {
	Enabled     bool   `env:"ENABLED" envDefault:"true"`
	ServiceURL  string `env:"SERVICE_URL" envDefault:"http://localhost:8080"`
	ServiceName string `env:"SERVICE_NAME" envDefault:"mandacode-ssam"`
	ClientName  string `env:"CLIENT_NAME" envDefault:"iam-client"`
}

func LoadManagementConfig() (*ManagementConfig, error) {
	cfg := &ManagementConfig{}
	if err := LoadConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}