package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Load loads configuration from environment variables into the provided struct
// T must be a pointer to a struct
func Load[T any](cfg T) error {
	if err := env.Parse(cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	return nil
}
