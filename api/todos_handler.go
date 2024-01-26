package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"

	"go-http/store"
)

type todoResponse struct {
	Id      uuid.UUID `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"director"`
}

type todoListResponse struct {
	Count int            `json:"count"`
	Data  []todoResponse `json:"data"`
}

type todoBody struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func NewTodoResponse(todo store.Todo) todoResponse {
	return todoResponse{
		Id:      todo.Id,
		Title:   todo.Title,
		Content: todo.Content,
	}
}

func NewTodoListResponse(todos []store.Todo) todoListResponse {
	var list []todoResponse

	for _, t := range todos {
		list = append(list, todoResponse{
			Id:      t.Id,
			Title:   t.Title,
			Content: t.Content,
		})
	}

	return todoListResponse{
		Count: len(list),
		Data:  list,
	}
}

// TODO: нужно разобраться зачем нужно Render и Bind дополнительно реализовывать
func (hr todoResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (hr todoListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) handleGetTodoList(w http.ResponseWriter, r *http.Request) {
	todos := s.store.GetAllTodos(r.Context())
	if todos == nil {
		render.Render(w, r, ErrNotFound)
		return
	}

	render.Render(w, r, NewTodoListResponse(todos))
}

func (s *Server) handleGetTodo(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	todo, err := s.store.GetTodo(r.Context(), id)
	if err != nil {
		render.Render(w, r, ErrInternalServerError)
		return
	}

	render.Render(w, r, NewTodoResponse(todo))
}

func (s *Server) handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.NewUUID()

	var body todoBody

	// TODO: изучить валидацию запросов
	err := render.DecodeJSON(r.Body, &body)
	if err != nil {
		fmt.Println(err.Error())
	}

	if body.Title == "" {
		render.Render(w, r, ErrBadRequest)
		return
	}

	todo := store.Todo{
		Id:      id,
		Title:   body.Title,
		Content: body.Content,
	}

	t, e := s.store.CreateTodo(r.Context(), todo)
	if e != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	render.Render(w, r, NewTodoResponse(t))
}

func (s *Server) handleUpdateTodo(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	var body todoBody

	// TODO: изучить валидацию запросов
	e := render.DecodeJSON(r.Body, &body)
	if e != nil {
		fmt.Println(err.Error())
	}

	if body.Title == "" {
		render.Render(w, r, ErrBadRequest)
		return
	}

	todo := store.Todo{
		Id:      id,
		Title:   body.Title,
		Content: body.Content,
	}

	t, err := s.store.UpdateTodo(r.Context(), todo)
	if err != nil {
		render.Render(w, r, ErrInternalServerError)
		return
	}

	render.Render(w, r, NewTodoResponse(t))
}

func (s *Server) handleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	e := s.store.DeleteTodo(r.Context(), id)
	if e != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
