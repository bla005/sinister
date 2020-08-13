package sinister

import (
	"fmt"

	"go.uber.org/zap"
)

// LogLevel ...
type LogLevel int

const (
	// DEBUG ...
	DEBUG LogLevel = iota
	// ERROR ...
	ERROR
	// INFO ...
	INFO
	// FATAL ...
	FATAL
	// WARN ...
	WARN
)

func newLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("failed wtf nigga")
		panic(err)
	}
	return logger
}

func newRouteLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("failed wtf nigga")
		panic(err)
	}
	return logger.With()
}
