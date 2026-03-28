package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"
)

type Configuration struct {
	serverAddress string `env:"ADDRESS"`
	logger        *zap.Logger
}

var cfg Configuration

func NewConfig() *Configuration {
	cfg = Configuration{}

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.serverAddress == "" {
		flag.StringVar(&cfg.serverAddress, "a", "localhost:8080", "address to listen on")
		flag.Parse()
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	cfg.logger = logger

	return &cfg
}

func Get() *Configuration {
	return &cfg
}

func (config *Configuration) GetServerAddress() string {
	return config.serverAddress
}

func (config *Configuration) GetLogger() *zap.SugaredLogger {
	return config.logger.Sugar()
}
