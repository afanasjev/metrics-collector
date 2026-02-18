package config

import (
	"go.uber.org/zap"
)

type Config interface {
	GetServerAddress() string
	GetLogger() *zap.SugaredLogger
}
