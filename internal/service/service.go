package service

import (
	"log/slog"
	"ppTodolistService/internal/repository"
)

type Service struct {
}

func MustNew(store repository.Repository, lg *slog.Logger) *Service {
	return nil
}
