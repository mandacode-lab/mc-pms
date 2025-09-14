package configs

import "bytes"

type RedisConfig struct {
	Addr             string `env:"REDIS_ADDR" envDefault:"localhost:6379" validate:"required,hostname_port"`
	Password         string `env:"REDIS_PASSWORD" envDefault:""`
	DB               int    `env:"REDIS_DB" envDefault:"0" validate:"min=0"`
	Mode             string `env:"REDIS_MODE" envDefault:"standard" validate:"oneof=standard sentinel cluster"`
	SentinelMaster   string `env:"REDIS_SENTINEL_MASTER" envDefault:"" validate:"required_if=Mode sentinel"`
	SentinelPassword string `env:"REDIS_SENTINEL_PASSWORD" envDefault:""`
}

func (cfg *RedisConfig) GetAddrList() []string {
	if cfg.Addr == "" {
		return nil
	}
	addrs := make([]string, 0)
	for addr := range bytes.SplitSeq([]byte(cfg.Addr), []byte(",")) {
		addrs = append(addrs, string(bytes.TrimSpace(addr)))
	}
	return addrs
}
