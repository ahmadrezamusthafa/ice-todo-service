package server

import (
	"github.com/ahmadrezamusthafa/ice-todo-service/adapter/api"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func (s *FiberServer) RegisterHandlers(todoHandler *api.TodoHandler, fileHandler *api.FileHandler) {
	s.app.Use(cors.New())
	s.app.Use(recover.New())
	s.app.Use(createLoggerMiddleware(s.logger))

	s.app.Get("/todo", todoHandler.GetAllTodos)
	s.app.Post("/todo", todoHandler.CreateTodo)
	s.app.Get("/todo/:id", todoHandler.GetTodo)
	s.app.Put("/todo/:id", todoHandler.UpdateTodo)
	s.app.Delete("/todo/:id", todoHandler.DeleteTodo)

	s.app.Post("/upload", fileHandler.UploadFile)
	s.app.Get("/file/:id", fileHandler.GetFile)
}
