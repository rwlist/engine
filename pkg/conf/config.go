package conf

import (
	"time"

	"github.com/caarlos0/env/v6"
)

type App struct {
	PrometheusBind     string        `env:"PROMETHEUS_BIND" envDefault:":2112"`
	DatabaseFile       string        `env:"DATABASE_FILE" envDefault:"data.db"`
	CompactionInterval time.Duration `env:"DATABASE_COMPACTION" envDefault:"1m"`

	ServerBind string `env:"SERVER_BIND" envDefault:":8080"`
}

func ParseEnv() (*App, error) {
	cfg := App{}
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
