package logger

import (
	"log/slog"
	"os"
)

var Logger = NewLogger()

func NewLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
