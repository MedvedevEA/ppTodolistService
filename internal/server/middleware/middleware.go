package middleware

import (
	"log/slog"
	srvErr "ppTodolistService/internal/server/err"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetLoggerMiddlewareFunc(lg *slog.Logger, appName string) func(c *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		start := time.Now()

		err := ctx.Next()

		lg.Info(
			"",
			slog.String("owner", "server"),
			slog.Any("method", ctx.Method()),
			slog.Any("path", ctx.Path()),
			slog.Any("statusCode", ctx.Response().StatusCode()),
			slog.Any("time", time.Since(start)),
		)
		return err
	}
}

func BadRequest(ctx *fiber.Ctx) error {
	return ctx.Status(404).SendString(srvErr.ErrRouteNotFound.Error())
}
