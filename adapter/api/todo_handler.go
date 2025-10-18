package api

import (
	"github.com/ahmadrezamusthafa/ice-todo-service/adapter/dto"
	"github.com/ahmadrezamusthafa/ice-todo-service/adapter/dto/mapper"
	"github.com/ahmadrezamusthafa/ice-todo-service/domain/apperrors"
	"github.com/ahmadrezamusthafa/ice-todo-service/domain/validator"
	"github.com/ahmadrezamusthafa/ice-todo-service/infrastructure/logger"
	"github.com/ahmadrezamusthafa/ice-todo-service/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TodoHandler struct {
	todoUseCase   usecase.TodoUseCaseInterface
	todoValidator *validator.TodoValidator
	logger        logger.Logger
	errorHandler  *ErrorHandler
}

func NewTodoHandler(todoUseCase usecase.TodoUseCaseInterface, logger logger.Logger) *TodoHandler {
	return &TodoHandler{
		todoUseCase:   todoUseCase,
		todoValidator: validator.NewTodoValidator(),
		logger:        logger,
		errorHandler:  NewErrorHandler(logger),
	}
}

func (h *TodoHandler) GetTodoUseCase() usecase.TodoUseCaseInterface {
	return h.todoUseCase
}

func (h *TodoHandler) GetLogger() logger.Logger {
	return h.logger
}

func (h *TodoHandler) CreateTodo(c *fiber.Ctx) error {
	var req dto.CreateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return h.errorHandler.Handle(c, apperrors.NewValidationError("Invalid request body", err))
	}

	if err := h.todoValidator.ValidateDescription(req.Description); err != nil {
		return h.errorHandler.Handle(c, apperrors.NewValidationError("Invalid description", err))
	}

	if err := h.todoValidator.ValidateDueDate(req.DueDate); err != nil {
		return h.errorHandler.Handle(c, apperrors.NewValidationError("Invalid due date", err))
	}

	todo := mapper.TodoRequestToEntity(&req)

	createdTodo, err := h.todoUseCase.CreateTodo(todo)
	if err != nil {
		return h.errorHandler.Handle(c, err)
	}

	response := mapper.TodoEntityToResponse(createdTodo)
	return dto.RespondWithJSON(c, fiber.StatusCreated, response)
}

func (h *TodoHandler) UpdateTodo(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return h.errorHandler.Handle(c, apperrors.NewValidationError("Invalid todo ID format", err))
	}

	existingTodo, err := h.todoUseCase.GetTodoByID(id)
	if err != nil {
		return h.errorHandler.Handle(c, err)
	}

	var req dto.UpdateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return h.errorHandler.Handle(c, apperrors.NewValidationError("Invalid request body", err))
	}

	if req.Description != nil {
		if err := h.todoValidator.ValidateDescription(*req.Description); err != nil {
			return h.errorHandler.Handle(c, apperrors.NewValidationError("Invalid description", err))
		}
	}

	if req.DueDate != nil {
		if err := h.todoValidator.ValidateDueDate(*req.DueDate); err != nil {
			return h.errorHandler.Handle(c, apperrors.NewValidationError("Invalid due date", err))
		}
	}

	existingTodo = mapper.UpdateTodoRequestToEntity(&req, existingTodo)

	updatedTodo, err := h.todoUseCase.UpdateTodo(id, existingTodo)
	if err != nil {
		return h.errorHandler.Handle(c, err)
	}

	response := mapper.TodoEntityToResponse(updatedTodo)
	return dto.RespondWithJSON(c, fiber.StatusOK, response)
}

func (h *TodoHandler) GetTodo(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return h.errorHandler.Handle(c, apperrors.NewValidationError("Invalid todo ID format", err))
	}

	todo, err := h.todoUseCase.GetTodoByID(id)
	if err != nil {
		return h.errorHandler.Handle(c, err)
	}

	response := mapper.TodoEntityToResponse(todo)
	return dto.RespondWithJSON(c, fiber.StatusOK, response)
}

func (h *TodoHandler) GetAllTodos(c *fiber.Ctx) error {
	todos, err := h.todoUseCase.GetAllTodos()
	if err != nil {
		return h.errorHandler.Handle(c, err)
	}

	response := mapper.TodoEntitiesToResponse(todos)
	return dto.RespondWithJSON(c, fiber.StatusOK, response)
}

func (h *TodoHandler) DeleteTodo(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return h.errorHandler.Handle(c, apperrors.NewValidationError("Invalid todo ID format", err))
	}

	err = h.todoUseCase.DeleteTodo(id)
	if err != nil {
		return h.errorHandler.Handle(c, err)
	}

	response := mapper.CreateDeleteResponse(id, "Todo item deleted successfully")
	return dto.RespondWithJSON(c, fiber.StatusOK, response)
}
