package config

import (
	"fmt"
	"time"

	"github.com/vrischmann/envconfig"
)

type Config struct {
	Producer
	Schedular
}

type Schedular struct {
	Period time.Duration
}

type Producer struct {
	Servers  string
	Protocol string
	Acks     string
}

func InitConfig(prefix string) (*Config, error) {
	conf := &Config{}
	if err := envconfig.InitWithPrefix(conf, prefix); err != nil {
		return nil, fmt.Errorf("init config error: %w", err)
	}

	return conf, nil
}
