package service

import (
	"log/slog"
	"ppTodolistService/internal/repository"
)

type Service struct {
	store repository.Repository
	lg    *slog.Logger
}

func MustNew(store repository.Repository, lg *slog.Logger) *Service {
	return &Service{
		store,
		lg,
	}
}
