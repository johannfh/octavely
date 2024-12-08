package main

import (
	"log/slog"
	"os"

	"github.com/johannfh/go-utils/assert"
	"github.com/johannfh/octavely/backend/api"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// TODO: get from cli flag --port (default "8080")
	addr := ":8080"
	server := api.Server{
		ListenAddr: addr,
		Logger:     logger,
	}

	slog.Info("server started", "addr", addr)
	err := server.Listen()
	assert.NoError(err, "failed to start server")

}
