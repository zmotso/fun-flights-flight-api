package app

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	FlyProviders []string
	Port         string
	JaegerURL    string
}

func NewConfig() (*Config, error) {
	conf := &Config{
		Port: ":8080",
	}

	conf.FlyProviders = strings.Split(os.Getenv("FLY_PROVIDERS"), ",")
	if os.Getenv("FLY_PROVIDERS") == "" || len(conf.FlyProviders) == 0 {
		return nil, errors.New("failed to get fly providers")
	}

	if port := os.Getenv("PORT"); port != "" {
		conf.Port = ":" + port
	}

	if jaeger := os.Getenv("JAEGER_URL"); jaeger != "" {
		conf.JaegerURL = jaeger
	}

	return conf, nil
}
