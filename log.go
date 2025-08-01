package main

import (
	"flag"
	"log/slog"
	"os"
)

var (
	logDebug bool
	logInfo  bool
)

func init() {
	flag.BoolVar(&logDebug, "vv", false, "Enable more verbose logging")
	flag.BoolVar(&logInfo, "v", false, "Enable more verbose logging")
}

func initializeLogger() {
	var logLevel slog.Level
	switch {
	case logDebug:
		logLevel = slog.LevelDebug
	case logInfo:
		logLevel = slog.LevelInfo
	default:
		logLevel = slog.LevelError
	}

	opts := slog.HandlerOptions{
		Level: logLevel,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &opts))
	slog.SetDefault(logger)
}
