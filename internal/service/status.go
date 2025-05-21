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

func (s *Service) AddStatus(ctx *fiber.Ctx) error {
	req := new(svcDto.AddStatus)
	if err := ctx.BodyParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.AddStatus"))
		return ctx.Status(400).SendString(svcErr.ErrBodyParse.Error())
	}
	if err := validator.Validate(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.AddStatus"))
		return ctx.Status(400).SendString(svcErr.ErrValidate.Error())
	}
	status, err := s.store.AddStatus(req.Name)
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(201).JSON(status)
}
func (s *Service) GetStatus(ctx *fiber.Ctx) error {
	req := new(svcDto.GetStatus)
	if err := ctx.ParamsParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.GetStatus"))
		return ctx.Status(400).SendString(svcErr.ErrParamsParse.Error())
	}
	if err := validator.Validate(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.GetStatus"))
		return ctx.Status(400).SendString(svcErr.ErrValidate.Error())
	}
	status, err := s.store.GetStatus(req.StatusId)
	if err != nil {
		if errors.Is(err, repoErr.ErrRecordNotFound) {
			return ctx.Status(404).SendString(err.Error())
		}
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(status)
}
func (s *Service) GetStatuses(ctx *fiber.Ctx) error {
	statuses, err := s.store.GetStatuses()
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(statuses)
}
func (s *Service) UpdateStatus(ctx *fiber.Ctx) error {
	req := new(svcDto.UpdateStatus)
	if err := ctx.BodyParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.UpdateStatus"))
		return ctx.Status(400).SendString(svcErr.ErrBodyParse.Error())
	}
	if err := ctx.ParamsParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.UpdateStatus"))
		return ctx.Status(400).SendString(svcErr.ErrParamsParse.Error())
	}
	if err := validator.Validate(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.UpdateStatus"))
		return ctx.Status(400).SendString(svcErr.ErrValidate.Error())
	}
	Status, err := s.store.UpdateStatus(&repoDto.UpdateStatus{
		StatusId: req.StatusId,
		Name:     req.Name,
	})
	if err != nil {
		if errors.Is(err, repoErr.ErrRecordNotFound) {
			return ctx.Status(404).SendString(err.Error())
		}
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(Status)
}
func (s *Service) RemoveStatus(ctx *fiber.Ctx) error {
	req := new(svcDto.RemoveStatus)
	if err := ctx.ParamsParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.RemoveStatus"))
		return ctx.Status(400).SendString(svcErr.ErrParamsParse.Error())
	}
	if err := validator.Validate(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.RemoveStatus"))
		return ctx.Status(400).SendString(svcErr.ErrValidate.Error())
	}
	err := s.store.RemoveStatus(req.StatusId)
	if err != nil {
		if errors.Is(err, repoErr.ErrRecordNotFound) {
			return ctx.Status(404).SendString(err.Error())
		}
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.SendStatus(204)
}
