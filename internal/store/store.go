package store

import (
	"log/slog"
	"ppTodolistService/internal/config"
)

type Store struct {
}

func MustNew(lg *slog.Logger, cfg *config.Store) *Store {
	return nil

}
