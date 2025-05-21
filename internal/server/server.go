package server

import (
	"log/slog"
	"os"
	"os/signal"
	"ppTodolistService/internal/config"
	"ppTodolistService/internal/server/middleware"
	"ppTodolistService/internal/service"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Server struct {
	app *fiber.App
	lg  *slog.Logger
	cfg *config.Server
}

func MustNew(svc *service.Service, lg *slog.Logger, cfg *config.Server) *Server {
	app := fiber.New(fiber.Config{
		AppName:      cfg.Name,
		WriteTimeout: cfg.WriteTimeout,
	})
	app.Use(recover.New(recover.ConfigDefault))
	app.Use(middleware.GetLoggerMiddlewareFunc(lg, cfg.Name))

	apiGroup := app.Group("/api")
	v1Group := apiGroup.Group("/v1")

	v1Group.Post("/logout", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
	v1Group.Post("/unregistration", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })

	v1Group.Post("/messages", svc.AddMessage)
	v1Group.Get("/messages/:messageId", svc.GetMessage)
	v1Group.Get("/messages", svc.GetMessages)
	v1Group.Patch("/messages/:messageId", svc.UpdateMessage)
	v1Group.Delete("/messages/:messageId", svc.RemoveMessage)

	v1Group.Post("/statuses", svc.AddStatus)
	v1Group.Get("/statuses/:statusId", svc.GetStatus)
	v1Group.Get("/statuses", svc.GetStatuses)
	v1Group.Patch("/statuses/:statusId", svc.UpdateStatus)
	v1Group.Delete("/statuses/:statusId", svc.RemoveStatus)

	v1Group.Post("/tasks", svc.AddTask)
	v1Group.Get("/tasks/:taskId", svc.GetTask)
	v1Group.Get("/tasks", svc.GetTasks)
	v1Group.Patch("/tasks/:taskId", svc.UpdateTask)
	v1Group.Delete("/tasks/:taskId", svc.RemoveTask)

	v1Group.Post("/taskusers", svc.AddTaskUser)
	v1Group.Get("/taskusers", svc.GetTaskUsers)
	v1Group.Delete("/taskusers/:taskUserId", svc.RemoveTaskUser)

	v1Group.Post("/users", svc.AddUser)
	v1Group.Get("/users", svc.GetUsers)
	v1Group.Delete("/users/:userId", svc.RemoveUser)

	app.Use(middleware.BadRequest)

	return &Server{
		app,
		lg,
		cfg,
	}
}

func (s *Server) Start() {
	chError := make(chan error, 1)
	go func() {
		s.lg.Info("server is started", slog.String("owner", "server"), slog.String("bindAddress", s.cfg.BindAddr))
		if err := s.app.Listen(s.cfg.BindAddr); err != nil {
			chError <- err
		}
	}()
	go func() {
		chQuit := make(chan os.Signal, 1)
		signal.Notify(chQuit, syscall.SIGINT, syscall.SIGTERM)
		<-chQuit
		chError <- s.app.Shutdown()
	}()
	if err := <-chError; err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "server"))
		return
	}
	s.lg.Info("server is stoped", slog.String("owner", "server"))

}
