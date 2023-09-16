package utils

import (
	"log/slog"
	"os"
)

func SetupLogger() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

}
