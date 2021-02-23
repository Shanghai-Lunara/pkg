package zaplogger

import (
	"go.uber.org/zap"
	"log"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Panic(err)
	}
}

func Sync() {
	_ = logger.Sync()
}

func Sugar() *zap.SugaredLogger {
	return logger.Sugar()
}

func Logger() *zap.Logger {
	return logger
}
