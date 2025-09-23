package shared

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

const SystemLevel zerolog.Level = 9

func init() {
	zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
		switch l {
		case SystemLevel:
			return "system"
		default:
			return l.String()
		}
	}
}

func SetupLogger(env string) zerolog.Logger {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	logger := zerolog.New(output).With().Timestamp().Logger()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if env == "dev" || env == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	return logger
}

func SystemLog(l zerolog.Logger) *zerolog.Event {
	return l.WithLevel(SystemLevel)
}
