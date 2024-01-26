package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Mkukarev/go-grpc/proto"
)

type Todo struct {
	Id      uuid.UUID
	Title   string
	Content string
}

type todoStore struct {
	grpcUrl string
	client  pb.TodoCRUDClient
}

func NewTodoStore(grpcUrl string) *todoStore {
	return &todoStore{
		grpcUrl: grpcUrl,
	}
}

func (store *todoStore) connect(ctx context.Context) {
	var opts = grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.Dial(store.grpcUrl, opts)
	if err != nil {
		fmt.Println("error connect to grpc")
		fmt.Println(err.Error())
	}

	client := pb.NewTodoCRUDClient(conn)

	store.client = client
}

func (store *todoStore) close(ctx context.Context) {
	// store.client.Close()
}

func (store *todoStore) CreateTodo(ctx context.Context, todo Todo) (Todo, error) {
	store.connect(ctx)
	defer store.close(ctx)

	args := &pb.CreateTodoMessage{
		Title:   todo.Title,
		Content: todo.Content,
	}

	t, err := store.client.CreateTodo(ctx, args)
	if err != nil {
		fmt.Println("unable to insert row:", err.Error())
		return Todo{}, err
	}

	id, _ := uuid.Parse(t.Id)

	responseTodo := Todo{
		Id:      id,
		Title:   t.Title,
		Content: t.Content,
	}

	return responseTodo, nil
}

func (store *todoStore) GetAllTodos(ctx context.Context) []Todo {
	store.connect(ctx)
	defer store.close(ctx)

	rows, err := store.client.GetAllTodos(ctx, &pb.Empty{})
	if err != nil {
		fmt.Println(123, err.Error())
		return nil
	}

	var todos []Todo

	for _, v := range rows.GetTodos() {
		id, _ := uuid.Parse(v.Id)
		todos = append(todos, Todo{Id: id, Title: v.Title, Content: v.Content})
	}

	return todos
}

func (store *todoStore) GetTodo(ctx context.Context, id uuid.UUID) (Todo, error) {
	store.connect(ctx)
	defer store.close(ctx)

	t, err := store.client.GetTodo(ctx, &pb.Id{Id: id.String()})

	if err != nil {
		return Todo{}, err
	}

	todo := Todo{
		Id:      id,
		Title:   t.Title,
		Content: t.Content,
	}

	return todo, nil
}

// TODO: изучить частичное изменение
func (store *todoStore) UpdateTodo(ctx context.Context, todo Todo) (Todo, error) {
	store.connect(ctx)
	defer store.close(ctx)

	t := &pb.Todo{
		Id:      todo.Id.String(),
		Title:   todo.Title,
		Content: todo.Content,
	}

	_, err := store.client.UpdateTodo(ctx, t)
	if err != nil {
		fmt.Println("delete row error:", err.Error())
		return Todo{}, err
	}

	return todo, nil
}

func (store *todoStore) DeleteTodo(ctx context.Context, id uuid.UUID) error {
	store.connect(ctx)
	defer store.close(ctx)

	_, err := store.client.DeleteTodo(ctx, &pb.Id{Id: id.String()})
	if err != nil {
		fmt.Println("delete row error:", err.Error())
		return err
	}

	return nil
}

type Interface interface {
	CreateTodo(ctx context.Context, todo Todo) (Todo, error)
	GetAllTodos(ctx context.Context) []Todo
	GetTodo(ctx context.Context, id uuid.UUID) (Todo, error)
	UpdateTodo(ctx context.Context, todo Todo) (Todo, error)
	DeleteTodo(ctx context.Context, id uuid.UUID) error
}
