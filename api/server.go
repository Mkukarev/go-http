package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"go-http/config"
	"go-http/store"
)

type Server struct {
	cfg    config.HTTPServer
	store  store.Interface
	router *chi.Mux
}

func NewServer(cfg config.HTTPServer, store store.Interface) *Server {
	s := &Server{
		cfg:    cfg,
		store:  store,
		router: chi.NewRouter(),
	}

	s.routes()

	return s
}

func (s *Server) Start(ctx context.Context) {
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Port),
		Handler:      s.router,
		IdleTimeout:  s.cfg.IdleTimeout,
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
	}

	server.ListenAndServe()
}
