package agent

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"
)

type Configuration struct {
	ServerAddress  string `env:"ADDRESS"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	Logger         *zap.Logger
}

var agentCfg Configuration

func NewConfig() *Configuration {
	agentCfg = Configuration{}

	err := env.Parse(&agentCfg)
	if err != nil {
		log.Fatal(err)
	}

	var serverAddress string
	var pollInterval int
	var reportInterval int
	flag.StringVar(&serverAddress, "a", "localhost:8080", "The server address in the format of host:port")
	flag.IntVar(&pollInterval, "p", 10, "The poll interval in seconds")
	flag.IntVar(&reportInterval, "r", 16, "The report interval in seconds")
	flag.Parse()

	if agentCfg.ServerAddress == "" {
		agentCfg.ServerAddress = serverAddress
	}
	agentCfg.ServerAddress = "http://" + agentCfg.ServerAddress

	if agentCfg.PollInterval == 0 {
		agentCfg.PollInterval = pollInterval
	}

	if agentCfg.ReportInterval == 0 {
		agentCfg.ReportInterval = reportInterval
	}

	agentCfg.Logger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer agentCfg.Logger.Sync()

	agentCfg.GetLogger().Infof("Configuration loaded: %+v", agentCfg)
	return &agentCfg
}

func Get() *Configuration {
	return &agentCfg
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
