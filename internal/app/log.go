package app

import (
	log "log/slog"
	"os"
)

func configureLogger() {
	handler := log.NewJSONHandler(os.Stdout, &log.HandlerOptions{
		Level: log.LevelDebug,
	})
	logger := log.New(handler)
	log.SetDefault(logger)
}
