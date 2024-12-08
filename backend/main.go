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
	port := 8080

	opts := append(
		api.DefaultOpts,
		api.WithHost("localhost"),
		api.WithPort(port),
	)
	server := api.NewServer(opts...)

	err := server.Listen()
	assert.NoError(err, "failed to start server")

}
