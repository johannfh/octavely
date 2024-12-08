package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	ListenAddr string
	Logger     *slog.Logger
}

func (s *Server) Listen() error {
	r := chi.NewRouter()

	err := http.ListenAndServe(s.ListenAddr, r)

	return err
}
