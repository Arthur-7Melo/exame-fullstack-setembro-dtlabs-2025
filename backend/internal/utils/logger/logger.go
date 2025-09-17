package logger

import (
	"log/slog"
	"os"
)

type CustomLogger struct {
    *slog.Logger
}

var Logger *CustomLogger

func init() {
    baseLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
    Logger = &CustomLogger{baseLogger}
    slog.SetDefault(baseLogger)
}