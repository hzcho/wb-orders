package config

import (
	"fmt"
	"time"

	"github.com/vrischmann/envconfig"
)

type Config struct {
	Consumer
	Server
	PG
	Cache
}

type Consumer struct {
	Brokers string
	Topics  []string
	GroupId string
	Offset  string
}

type Server struct {
	Port      string
	ReadTime  time.Duration
	WriteTime time.Duration
}

type PG struct {
	Username string
	Host     string
	Port     string
	DBName   string
	Password string
}

type Cache struct {
	Capacity int
}

func InitConfig(prefix string) (*Config, error) {
	conf := &Config{}
	if err := envconfig.InitWithPrefix(conf, prefix); err != nil {
		return nil, fmt.Errorf("init config error: %w", err)
	}

	return conf, nil
}
