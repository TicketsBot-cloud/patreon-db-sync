package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/google/uuid"
	"go.uber.org/zap/zapcore"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Daemon           bool          `env:"DAEMON" envDefault:"true"`
	RunFrequency     time.Duration `env:"RUN_FREQUENCY" envDefault:"1m"`
	ExecutionTimeout time.Duration `env:"EXECUTION_TIMEOUT" envDefault:"3m"`

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

	TierSkus map[uint64]uuid.UUID `env:"TIER_SKUS"`
}

func LoadFromEnv() (Config, error) {
	var config Config
	err := env.ParseWithOptions(&config, env.Options{
		FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeOf(map[uint64]uuid.UUID{}): tierMapParser,
		},
	})
	return config, err
}

func tierMapParser(value string) (interface{}, error) {
	m := make(map[uint64]uuid.UUID)
	values := strings.Split(value, ",")

	for _, pair := range values {
		split := strings.Split(pair, "=")
		if len(split) != 2 {
			return nil, fmt.Errorf("invalid format: %s", pair)
		}

		key, err := strconv.ParseUint(split[0], 10, 64)
		if err != nil {
			return nil, err
		}

		uuid, err := uuid.Parse(split[1])
		if err != nil {
			return nil, err
		}

		m[key] = uuid
	}

	return m, nil
}
