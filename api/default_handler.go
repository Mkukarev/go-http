package api

import (
	"net/http"

	"github.com/go-chi/render"
)

type defaultResponse struct {
	OK bool `json:"ok"`
}

func (hr defaultResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// TODO: переделать на текст "Todo list"
func (s *Server) handleGetDefault(w http.ResponseWriter, r *http.Request) {
	health := healthResponse{OK: true}
	render.Render(w, r, health)
}
