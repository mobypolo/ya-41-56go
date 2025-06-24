package logger

import (
	"go.uber.org/zap"
	"log"
	"sync"
)

var (
	loggerInstance *zap.Logger
	once           sync.Once
	isLoadedLogger bool
)

func Init(mode ModeLogger) {
	once.Do(func() {
		var err error
		if mode.IsDev() {
			loggerInstance, err = zap.NewDevelopment()
		} else {
			loggerInstance, err = zap.NewProduction()
		}
		if err != nil {
			log.Fatalf("failed to init zap logger: %v", err)
		}
		isLoadedLogger = true
	})
}

func L() *zap.Logger {
	if !isLoadedLogger {
		log.Fatal("logger not initialized — call logger.Init() first")
	}
	return loggerInstance
}
