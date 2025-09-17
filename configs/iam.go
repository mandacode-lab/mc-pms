package configs

type IAMConfig struct {
	Enabled     bool   `env:"ENABLED" envDefault:"true"`
	ServiceURL  string `env:"SERVICE_URL" envDefault:"http://localhost:8080"`
}
