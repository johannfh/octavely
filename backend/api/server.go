package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/johannfh/go-utils/assert"
	"github.com/johannfh/octavely/backend/db/repository"
)

type Server struct {
	host string
	port int

	logger  *slog.Logger
	queries *repository.Queries
}

type NewServerOpt func(*Server)

func NewServer(options ...NewServerOpt) *Server {
	server := &Server{}

	for _, option := range options {
		option(server)
	}

	assert.NotNil(server.queries, "database queries missing")

	return server
}

var DefaultOpts = []NewServerOpt{
	WithHost("localhost"),
	WithPort(8080),
	WithLogger(slog.Default()),
}

func WithHost(host string) NewServerOpt {
	return func(s *Server) {
		s.host = host
	}
}

func WithPort(port int) NewServerOpt {
	return func(s *Server) {
		s.port = port
	}
}

func WithLogger(logger *slog.Logger) NewServerOpt {
	return func(s *Server) {
		s.logger = logger
	}
}

func WithQueries(queries *repository.Queries) NewServerOpt {
	return func(s *Server) {
		s.queries = queries
	}
}

func (s *Server) Listen() (err error) {
	r := chi.NewRouter()

	r.Get("/composers/{id}", s.handleGetComposer)

	r.Get("/composers", s.handleGetAllComposers)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	slog.Info("server listening", "host", s.host, "port", s.port)
	err = http.ListenAndServe(addr, r)

	return err
}
