package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	ListenAddr string
}

func (s *Server) Listen() error {
	r := chi.NewRouter()

	http.ListenAndServe(s.ListenAddr, r)
	return nil
}
