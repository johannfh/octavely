package main

import (
	"database/sql"
	"errors"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/johannfh/go-utils/assert"
	"github.com/johannfh/octavely/backend/api"
	"github.com/johannfh/octavely/backend/db"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	dbPath := "data/db.sqlite3"

	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		slog.Error("failed to create database connection", "err", err)
		os.Exit(1)
	}
	slog.Info("created database connection")

	source, err := iofs.New(db.Schema, "schema")
	if err != nil {
		slog.Error("failed to create migrator source instance", "err", err)
		os.Exit(1)
	}

	database, err := sqlite3.WithInstance(sqlDB, &sqlite3.Config{})
	if err != nil {
		slog.Error("failed to create migrator database instance", "err", err)
		os.Exit(1)
	}

	migrator, err := migrate.NewWithInstance("iofs", source, "sqlite3", database)
	if err != nil {
		slog.Error("failed to create migrator", "err", err)
		os.Exit(1)
	}

	oldVersion, _, _ := migrator.Version()

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("failed to migrate database", "err", err)
		os.Exit(1)
	}

	newVersion, _, _ := migrator.Version()

	slog.Info(
		"applied database migrations",
		"oldVersion", oldVersion,
		"newVersion", newVersion,
		"change", newVersion-oldVersion,
	)

	// TODO: get from cli flag --port (default "8080")
	port := 8080

	opts := append(
		api.DefaultOpts,
		api.WithHost("localhost"),
		api.WithPort(port),
	)
	server := api.NewServer(opts...)

	err = server.Listen()
	assert.NoError(err, "failed to start server")

}
