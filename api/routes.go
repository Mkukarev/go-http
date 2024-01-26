package api

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func (s *Server) routes() {
	s.router.Use(render.SetContentType(render.ContentTypeJSON))

	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Timeout(60 * time.Second))

	s.router.Get("/", s.handleGetDefault)

	s.router.Get("/health", s.handleGetHealth)

	s.router.Route("/api/todos", func(r chi.Router) {
		r.Get("/", s.handleGetTodoList)
		r.Post("/", s.handleCreateTodo)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", s.handleGetTodo)
			r.Put("/", s.handleUpdateTodo)
			r.Delete("/", s.handleDeleteTodo)
		})
	})
}
