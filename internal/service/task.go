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

func (s *Service) AddTask(ctx *fiber.Ctx) error {
	req := new(svcDto.AddTask)
	if err := ctx.BodyParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.AddTask"))
		return ctx.Status(400).SendString(svcErr.ErrBodyParse.Error())
	}
	if err := validator.Validate(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.AddTask"))
		return ctx.Status(400).SendString(svcErr.ErrValidate.Error())
	}
	task, err := s.store.AddTask(&repoDto.AddTask{
		StatusId:    req.StatusId,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(201).JSON(task)
}
func (s *Service) GetTask(ctx *fiber.Ctx) error {
	req := new(svcDto.GetTask)
	if err := ctx.ParamsParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.GetTask"))
		return ctx.Status(400).SendString(svcErr.ErrParamsParse.Error())
	}
	if err := validator.Validate(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.GetTask"))
		return ctx.Status(400).SendString(svcErr.ErrValidate.Error())
	}
	task, err := s.store.GetTask(req.TaskId)
	if err != nil {
		if errors.Is(err, repoErr.ErrRecordNotFound) {
			return ctx.Status(404).SendString(err.Error())
		}
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(task)
}
func (s *Service) GetTasks(ctx *fiber.Ctx) error {
	req := &svcDto.GetTasks{
		Offset: 0,
		Limit:  10,
	}
	if err := ctx.QueryParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.GetTasks"))
		return ctx.Status(400).SendString(svcErr.ErrQueryParse.Error())
	}
	if err := validator.Validate(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.GetTasks"))
		return ctx.Status(400).SendString(svcErr.ErrValidate.Error())
	}
	tasks, err := s.store.GetTasks(&repoDto.GetTasks{
		Offset:   req.Offset,
		Limit:    req.Limit,
		StatusId: req.StatusId,
	})
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(tasks)
}
func (s *Service) UpdateTask(ctx *fiber.Ctx) error {
	req := new(svcDto.UpdateTask)
	if err := ctx.BodyParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.UpdateTask"))
		return ctx.Status(400).SendString(svcErr.ErrBodyParse.Error())
	}
	if err := ctx.ParamsParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.UpdateTask"))
		return ctx.Status(400).SendString(svcErr.ErrParamsParse.Error())
	}
	if err := validator.Validate(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.UpdateTask"))
		return ctx.Status(400).SendString(svcErr.ErrValidate.Error())
	}
	task, err := s.store.UpdateTask(&repoDto.UpdateTask{
		TaskId:      req.TaskId,
		StatusId:    req.StatusId,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		if errors.Is(err, repoErr.ErrRecordNotFound) {
			return ctx.Status(404).SendString(err.Error())
		}
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(task)
}
func (s *Service) RemoveTask(ctx *fiber.Ctx) error {
	req := new(svcDto.RemoveTask)
	if err := ctx.ParamsParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.RemoveTask"))
		return ctx.Status(400).SendString(svcErr.ErrParamsParse.Error())
	}
	if err := validator.Validate(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.RemoveTask"))
		return ctx.Status(400).SendString(svcErr.ErrValidate.Error())
	}
	err := s.store.RemoveTask(req.TaskId)
	if err != nil {
		if errors.Is(err, repoErr.ErrRecordNotFound) {
			return ctx.Status(404).SendString(err.Error())
		}
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.SendStatus(204)
}
