package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	host string
	port int

	logger *slog.Logger
}

type NewServerOpt func(*Server)

func NewServer(options ...NewServerOpt) *Server {
	server := &Server{}

	for _, option := range options {
		option(server)
	}

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

func (s *Server) Listen() (err error) {
	r := chi.NewRouter()

	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	slog.Info("server listening", "host", s.host, "port", s.port)
	err = http.ListenAndServe(addr, r)

	return err
}
