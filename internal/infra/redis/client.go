package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr             string
	Password         string
	DB               int
	Mode             string
	SentinelMaster   string
	SentinelPassword string
}

func NewClient(ctx context.Context, cfg Config) (*redis.UniversalClient, error) {
	var client redis.UniversalClient

	switch cfg.Mode {
	case "standalone":
		client = redis.NewClient(&redis.Options{
			Addr:     cfg.Addr,
			Password: cfg.Password,
			DB:       cfg.DB,
		})
	case "sentinel":
		client = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:       cfg.SentinelMaster,
			SentinelAddrs:    []string{cfg.Addr},
			Password:         cfg.Password,
			SentinelPassword: cfg.SentinelPassword,
			DB:               cfg.DB,
		})
	case "cluster":
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    []string{cfg.Addr},
			Password: cfg.Password,
		})
	default:
		return nil, fmt.Errorf("unsupported redis mode: %s", cfg.Mode)
	}

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &client, nil
}