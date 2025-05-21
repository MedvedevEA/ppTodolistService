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

func (s *Service) AddUser(ctx *fiber.Ctx) error {
	req := new(svcDto.AddUser)
	if err := ctx.BodyParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.AddUser"))
		return ctx.Status(400).SendString(svcErr.ErrBodyParse.Error())
	}
	if err := validator.Validate(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.AddUser"))
		return ctx.Status(400).SendString(svcErr.ErrValidate.Error())
	}
	user, err := s.store.AddUserWithUserId(&repoDto.AddUser{
		UserId: req.UserId,
		Name:   req.Name,
	})
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(201).JSON(user)
}
func (s *Service) GetUsers(ctx *fiber.Ctx) error {
	req := &svcDto.GetUsers{
		Offset: 0,
		Limit:  10,
	}
	if err := ctx.QueryParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.GetUsers"))
		return ctx.Status(400).SendString(svcErr.ErrQueryParse.Error())
	}
	if err := validator.Validate(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.GetUsers"))
		return ctx.Status(400).SendString(svcErr.ErrValidate.Error())
	}
	users, err := s.store.GetUsers(&repoDto.GetUsers{
		Offset: req.Offset,
		Limit:  req.Limit,
		Name:   req.Name,
	})
	if err != nil {
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.Status(200).JSON(users)
}
func (s *Service) RemoveUser(ctx *fiber.Ctx) error {
	req := new(svcDto.RemoveUser)
	if err := ctx.ParamsParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.RemoveUser"))
		return ctx.Status(400).SendString(svcErr.ErrParamsParse.Error())
	}
	if err := validator.Validate(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.RemoveUser"))
		return ctx.Status(400).SendString(svcErr.ErrValidate.Error())
	}
	err := s.store.RemoveUser(req.UserId)
	if err != nil {
		if errors.Is(err, repoErr.ErrRecordNotFound) {
			return ctx.Status(404).SendString(err.Error())
		}
		return ctx.Status(500).SendString(err.Error())
	}
	return ctx.SendStatus(204)
}
