package service

import (
	"errors"
	"fmt"
	"log/slog"
	repoDto "ppTodolistService/internal/repository/dto"
	repoErr "ppTodolistService/internal/repository/err"
	svcDto "ppTodolistService/internal/service/dto"
	svcErr "ppTodolistService/internal/service/err"
	"ppTodolistService/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

func (s *Service) AddMessage(ctx *fiber.Ctx) error {
	req := new(svcDto.AddMessage)
	if err := ctx.BodyParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.AddMessage"))
		return ctx.Status(400).SendString(svcErr.ErrBodyParse.Error())
	}
	if err := validator.Validate(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.AddMessage"))
		return ctx.Status(400).SendString(svcErr.ErrValidate.Error())
	}
	message, err := s.store.AddMessage(&repoDto.AddMessage{
		TaskId: req.TaskId,
		UserId: req.UserId,
		Text:   req.Text,
	})
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(201).JSON(message)
}
func (s *Service) GetMessage(ctx *fiber.Ctx) error {
	req := new(svcDto.GetMessage)
	if err := ctx.ParamsParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.GetMessage"))
		return ctx.Status(400).SendString(svcErr.ErrParamsParse.Error())
	}
	if err := validator.Validate(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.GetMessage"))
		return ctx.Status(400).SendString(svcErr.ErrValidate.Error())
	}
	message, err := s.store.GetMessage(req.MessageId)
	if err != nil {
		if errors.Is(err, repoErr.ErrRecordNotFound) {
			return ctx.Status(404).SendString(err.Error())
		}
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(message)
}
func (s *Service) GetMessages(ctx *fiber.Ctx) error {
	req := &svcDto.GetMessages{
		Offset: 0,
		Limit:  10,
	}
	if err := ctx.QueryParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.GetMessages"))
		return ctx.Status(400).SendString(svcErr.ErrQueryParse.Error())
	}
	if err := validator.Validate(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.GetMessages"))
		return ctx.Status(400).SendString(svcErr.ErrValidate.Error())
	}
	Messages, err := s.store.GetMessages(&repoDto.GetMessages{
		Offset: req.Offset,
		Limit:  req.Limit,
		TaskId: req.TaskId,
	})
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(Messages)
}
func (s *Service) UpdateMessage(ctx *fiber.Ctx) error {
	req := new(svcDto.UpdateMessage)
	if err := ctx.BodyParser(req); err != nil {
		err := fmt.Errorf("service.UpdateMessage: %w (%v)", svcErr.ErrBodyParse, err)
		return ctx.Status(400).SendString(err.Error())
	}
	if err := ctx.ParamsParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.UpdateMessage"))
		return ctx.Status(400).SendString(svcErr.ErrParamsParse.Error())
	}
	if err := validator.Validate(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.UpdateMessage"))
		return ctx.Status(400).SendString(svcErr.ErrValidate.Error())
	}
	message, err := s.store.UpdateMessage(&repoDto.UpdateMessage{
		MessageId: req.MessageId,
		Text:      req.Text,
	})
	if err != nil {
		if errors.Is(err, repoErr.ErrRecordNotFound) {
			return ctx.Status(404).SendString(err.Error())
		}
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(message)
}
func (s *Service) RemoveMessage(ctx *fiber.Ctx) error {
	req := new(svcDto.RemoveMessage)
	if err := ctx.ParamsParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.RemoveMessage"))
		return ctx.Status(400).SendString(svcErr.ErrParamsParse.Error())
	}
	if err := validator.Validate(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.RemoveMessage"))
		return ctx.Status(400).SendString(svcErr.ErrValidate.Error())
	}
	err := s.store.RemoveMessage(req.MessageId)
	if err != nil {
		if errors.Is(err, repoErr.ErrRecordNotFound) {
			return ctx.Status(404).SendString(err.Error())
		}
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.SendStatus(204)
}
