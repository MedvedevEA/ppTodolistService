package main

import (
	"ppTodolistService/internal/config"
	"ppTodolistService/internal/logger"
	"ppTodolistService/internal/server"
	"ppTodolistService/internal/service"
	"ppTodolistService/internal/store"
)

func main() {

	cfg := config.MustNew()
	lg := logger.MustNew(cfg.Environment)
	store := store.MustNew(lg, &cfg.Store)
	service := service.MustNew(store, lg)
	server := server.MustNew(service, lg, &cfg.Server)

	server.Start()

}
