package api

import (
	"errors"
	"github.com/ahmadrezamusthafa/ice-todo-service/adapter/dto"
	"github.com/ahmadrezamusthafa/ice-todo-service/domain/apperrors"
	"github.com/ahmadrezamusthafa/ice-todo-service/infrastructure/logger"
	"github.com/gofiber/fiber/v2"
)

type ErrorHandler struct {
	logger logger.Logger
}

func NewErrorHandler(logger logger.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

func (h *ErrorHandler) Handle(c *fiber.Ctx, err error) error {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		switch appErr.Type {
		case apperrors.ErrorTypeInternal, apperrors.ErrorTypeDatabase, apperrors.ErrorTypeStorage:
			h.logger.Error("Error: %v", appErr)
		default:
			h.logger.Info("Error: %v", appErr)
		}

		return dto.RespondWithError(c, appErr.StatusCode(), appErr.Message)
	}

	h.logger.Error("Unexpected error: %v", err)
	return dto.RespondWithError(c, fiber.StatusInternalServerError, "An unexpected error occurred")
}

func (h *ErrorHandler) HandleErrorWithMessage(c *fiber.Ctx, err error, message string) error {
	h.logger.Error("%s: %v", message, err)

	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		return dto.RespondWithError(c, appErr.StatusCode(), message)
	}

	return dto.RespondWithError(c, fiber.StatusInternalServerError, message)
}
