package config

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const configFile = "data/config.yaml"

type Config struct {
	Token                    string        `yaml:"token"`
	GracefullShutdownTimeout time.Duration `yaml:"gracefull_shutdown_timeout"`
	UpdateCurrenciesInterval time.Duration `yaml:"update_currencies_interval"`
	CacheHost                string        `yaml:"cache_host"`
	CachePort                int           `yaml:"cache_port"`
	TopicReport              string        `yaml:"topic_report"`
	KafkaBrokers             []string      `yaml:"kafka_brokers"`
}

func New() (*Config, error) {
	c := &Config{}

	rawYAML, err := os.ReadFile(configFile)
	if err != nil {
		return nil, errors.Wrap(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, &c)
	if err != nil {
		return nil, errors.Wrap(err, "parsing yaml")
	}

	return c, nil
}
