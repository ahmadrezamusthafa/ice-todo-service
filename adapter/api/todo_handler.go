package api

import (
	"github.com/ahmadrezamusthafa/ice-todo-service/adapter/dto"
	"github.com/ahmadrezamusthafa/ice-todo-service/domain/entity"
	"github.com/ahmadrezamusthafa/ice-todo-service/domain/validator"
	"github.com/ahmadrezamusthafa/ice-todo-service/infrastructure/logger"
	"github.com/ahmadrezamusthafa/ice-todo-service/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TodoHandler struct {
	todoUseCase   *usecase.TodoUseCase
	todoValidator *validator.TodoValidator
	logger        logger.Logger
}

func NewTodoHandler(todoUseCase *usecase.TodoUseCase, logger logger.Logger) *TodoHandler {
	return &TodoHandler{
		todoUseCase:   todoUseCase,
		todoValidator: validator.NewTodoValidator(),
		logger:        logger,
	}
}

func (h *TodoHandler) GetTodoUseCase() *usecase.TodoUseCase {
	return h.todoUseCase
}

func (h *TodoHandler) GetLogger() logger.Logger {
	return h.logger
}

func (h *TodoHandler) CreateTodo(c *fiber.Ctx) error {
	var req dto.CreateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse request body: %v", err)
		return dto.RespondWithError(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	if err := h.todoValidator.ValidateDescription(req.Description); err != nil {
		return dto.RespondWithError(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.todoValidator.ValidateDueDate(req.DueDate); err != nil {
		return dto.RespondWithError(c, fiber.StatusBadRequest, err.Error())
	}

	todo := entity.NewTodoItem(req.Description, req.DueDate, req.FileID)

	createdTodo, err := h.todoUseCase.CreateTodo(todo)
	if err != nil {
		h.logger.Error("Failed to create todo item: %v", err)
		return dto.RespondWithError(c, fiber.StatusInternalServerError, "Failed to create todo item: internal server error")
	}

	return dto.RespondWithJSON(c, fiber.StatusCreated, createdTodo)
}

func (h *TodoHandler) UpdateTodo(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return dto.RespondWithError(c, fiber.StatusBadRequest, "Invalid todo ID format")
	}

	existingTodo, err := h.todoUseCase.GetTodoByID(id)
	if err != nil {
		h.logger.Error("Failed to get todo item: %v", err)
		return dto.RespondWithError(c, fiber.StatusNotFound, "Todo item not found")
	}

	var req dto.UpdateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse request body: %v", err)
		return dto.RespondWithError(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	if req.Description != nil {
		if err := h.todoValidator.ValidateDescription(*req.Description); err != nil {
			return dto.RespondWithError(c, fiber.StatusBadRequest, err.Error())
		}
		existingTodo.Description = *req.Description
	}

	if req.DueDate != nil {
		if err := h.todoValidator.ValidateDueDate(*req.DueDate); err != nil {
			return dto.RespondWithError(c, fiber.StatusBadRequest, err.Error())
		}
		existingTodo.DueDate = *req.DueDate
	}

	updatedTodo, err := h.todoUseCase.UpdateTodo(id, existingTodo)
	if err != nil {
		h.logger.Error("Failed to update todo item: %v", err)
		return dto.RespondWithError(c, fiber.StatusInternalServerError, "Failed to update todo item: internal server error")
	}

	return dto.RespondWithJSON(c, fiber.StatusOK, updatedTodo)
}
