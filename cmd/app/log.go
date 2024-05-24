package main

import (
	"log/slog"
	"os"
)

const logLevel = slog.LevelError

func init() {
	opts := slog.HandlerOptions{
		Level: logLevel,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &opts))
	slog.SetDefault(logger)
}
