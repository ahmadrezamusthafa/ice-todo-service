package server

import (
	"context"
	"fmt"
	"github.com/ahmadrezamusthafa/ice-todo-service/adapter/api"
	"github.com/ahmadrezamusthafa/ice-todo-service/adapter/dto"
	"github.com/ahmadrezamusthafa/ice-todo-service/config"
	"github.com/ahmadrezamusthafa/ice-todo-service/infrastructure/logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type FiberServer struct {
	config *config.Config
	logger logger.Logger
	app    *fiber.App
}

func NewFiberServer(config *config.Config, logger logger.Logger) *FiberServer {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return dto.RespondWithError(c, code, err.Error())
		},
	})

	return &FiberServer{
		config: config,
		logger: logger,
		app:    app,
	}
}

func (s *FiberServer) RegisterHandlers(todoHandler *api.TodoHandler, fileHandler *api.FileHandler) {
	s.app.Use(cors.New())
	s.app.Use(recover.New())
	s.app.Use(createLoggerMiddleware(s.logger))

	s.app.Post("/todo", todoHandler.CreateTodo)
	s.app.Post("/upload", fileHandler.UploadFile)
}

func (s *FiberServer) Start() error {
	go func() {
		s.logger.Info("Server starting on port %s", s.config.Server.Port)
		if err := s.app.Listen(fmt.Sprintf(":%s", s.config.Server.Port)); err != nil {
			s.logger.Fatal("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.logger.Info("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.app.ShutdownWithContext(ctx); err != nil {
		s.logger.Error("Server forced to shutdown: %v", err)
		return err
	}

	s.logger.Info("Server exited properly")
	return nil
}

func createLoggerMiddleware(log logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		log.Info("%s %s - %d - %v", c.Method(), c.Path(), c.Response().StatusCode(), time.Since(start))

		return err
	}
}
