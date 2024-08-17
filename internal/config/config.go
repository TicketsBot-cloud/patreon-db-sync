package config

import (
	env "github.com/caarlos0/env/v11"
	"go.uber.org/zap/zapcore"
	"time"
)
import "net/url"

type Config struct {
	Daemon       bool          `env:"DAEMON" envDefault:"true"`
	RunFrequency time.Duration `env:"RUN_FREQUENCY" envDefault:"1m"`

	SentryDsn string        `env:"SENTRY_DSN"`
	JsonLogs  bool          `env:"JSON_LOGS" envDefault:"false"`
	LogLevel  zapcore.Level `env:"LOG_LEVEL" envDefault:"info"`

	PatreonProxy struct {
		RootUrl   *url.URL `env:"ROOT_URL" envDefault:"http://localhost:8081"`
		AuthToken string   `env:"AUTH_TOKEN"`
	} `envPrefix:"PATREON_PROXY_"`

	DatabaseUri string `env:"DATABASE_URI"`

	MinEntitlementsThreshold int `env:"MIN_ENTITLEMENTS_THRESHOLD"`
	MaxRemovalsThreshold     int `env:"MAX_REMOVALS_THRESHOLD"`
	GracePeriodDays          int `env:"GRACE_PERIOD_DAYS" envDefault:"7"`
}

func LoadFromEnv() (Config, error) {
	var config Config
	err := env.Parse(&config)
	return config, err
}
