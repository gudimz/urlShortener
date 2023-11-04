package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Log struct {
	*zap.Logger
}

func New(cfg Config) *Log {
	var (
		log *zap.Logger
		err error
	)

	switch cfg.LogLevel {
	case "debug":
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		log, err = config.Build()
	default:
		config := zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		log, err = config.Build()
	}

	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(log)
	defer log.Sync()

	return &Log{log}
}
