package main

import (
	"context"
	"fmt"

	"go-http/api"
	"go-http/config"
	"go-http/store"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		fmt.Println(err.Error())
	}

	store := store.NewTodoStore(cfg.TodoGRPCServerUrl)

	server := api.NewServer(cfg.HTTPServer, store)
	server.Start(ctx)
}
