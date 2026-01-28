package server

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"
)

type Configuration struct {
	ServerAddress string `env:"ADDRESS"`
	Logger        *zap.Logger
	Debug         bool `env:"DEBUG"`
}

var srvCfg Configuration

func NewConfig() *Configuration {
	srvCfg = Configuration{}

	err := env.Parse(&srvCfg)
	if err != nil {
		log.Fatal(err)
	}

	var serverAddress string
	var debug bool
	flag.BoolVar(&debug, "debug", true, "debug flag") //Debug level
	flag.StringVar(&serverAddress, "a", ":8080", "address to listen on")
	flag.Parse()

	if srvCfg.ServerAddress == "" {
		srvCfg.ServerAddress = serverAddress
	}

	if !srvCfg.Debug {
		srvCfg.Debug = debug
	}

	srvCfg.Logger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer srvCfg.Logger.Sync()

	return &srvCfg
}

func GetServerConfig() *Configuration {
	return &srvCfg
}

func (config *Configuration) GetServerAddress() string {
	return config.ServerAddress
}

func (config *Configuration) GetLogger() *zap.SugaredLogger {
	if config == nil || config.Logger == nil {
		return zap.NewNop().Sugar()
	}
	return config.Logger.Sugar()
}
