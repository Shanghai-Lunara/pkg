package zaplogger

import (
	"go.uber.org/zap"
	"log"
)

var logger *zap.Logger
var sugar *zap.SugaredLogger

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Panic(err)
	}
	sugar = logger.Sugar()
}

func Sync() {
	_ = logger.Sync()
}

func Sugar() *zap.SugaredLogger {
	return sugar
}

func Logger() *zap.Logger {
	return logger
}
