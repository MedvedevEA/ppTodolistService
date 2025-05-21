package service

import (
	"errors"
	"log/slog"
	repoDto "ppTodolistService/internal/repository/dto"
	repoErr "ppTodolistService/internal/repository/err"
	svcDto "ppTodolistService/internal/service/dto"
	svcErr "ppTodolistService/internal/service/err"
	"ppTodolistService/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

func (s *Service) AddTaskUser(ctx *fiber.Ctx) error {
	req := new(svcDto.AddTaskUser)
	if err := ctx.BodyParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.AddTaskUser"))
		return ctx.Status(400).SendString(svcErr.ErrBodyParse.Error())
	}
	if err := validator.Validate(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.AddTaskUser"))
		return ctx.Status(400).SendString(svcErr.ErrValidate.Error())
	}
	taskUser, err := s.store.AddTaskUser(&repoDto.AddTaskUser{
		TaskId: req.TaskId,
		UserId: req.UserId,
	})
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(201).JSON(taskUser)
}
func (s *Service) GetTaskUsers(ctx *fiber.Ctx) error {
	req := &svcDto.GetTaskUsers{
		Offset: 0,
		Limit:  10,
	}
	if err := ctx.QueryParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.GetTaskUsers"))
		return ctx.Status(400).SendString(svcErr.ErrQueryParse.Error())
	}
	if err := validator.Validate(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.GetTaskUsers"))
		return ctx.Status(400).SendString(svcErr.ErrValidate.Error())
	}
	taskUsers, err := s.store.GetTaskUsers(&repoDto.GetTaskUsers{
		Offset: req.Offset,
		Limit:  req.Limit,
		TaskId: req.TaskId,
		UserId: req.UserId,
	})
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(taskUsers)
}
func (s *Service) RemoveTaskUser(ctx *fiber.Ctx) error {
	req := new(svcDto.RemoveTaskUser)
	if err := ctx.ParamsParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.RemoveTaskUser"))
		return ctx.Status(400).SendString(svcErr.ErrParamsParse.Error())
	}
	if err := validator.Validate(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.RemoveTaskUser"))
		return ctx.Status(400).SendString(svcErr.ErrValidate.Error())
	}
	err := s.store.RemoveTaskUser(req.TaskUserId)
	if err != nil {
		if errors.Is(err, repoErr.ErrRecordNotFound) {
			return ctx.Status(404).SendString(err.Error())
		}
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.SendStatus(204)
}
