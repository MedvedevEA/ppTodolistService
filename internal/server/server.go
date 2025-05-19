package server

import (
	"log/slog"
	"ppTodolistService/internal/config"
	"ppTodolistService/internal/service"
)

type Server struct {
}

func MustNew(service *service.Service, lg *slog.Logger, cfg *config.Server) *Server {
	return nil
}

func (s *Server) Start() {

}

func (s *Server) Stop() {

}
