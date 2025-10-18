package api

import (
	"github.com/ahmadrezamusthafa/ice-todo-service/adapter/dto"
	"github.com/ahmadrezamusthafa/ice-todo-service/domain/apperrors"
	"github.com/ahmadrezamusthafa/ice-todo-service/infrastructure/logger"
	"github.com/ahmadrezamusthafa/ice-todo-service/usecase"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type StreamHandler struct {
	streamUseCase usecase.StreamUseCaseInterface
	logger        logger.Logger
	errorHandler  *ErrorHandler
}

func NewStreamHandler(streamUseCase usecase.StreamUseCaseInterface, logger logger.Logger) *StreamHandler {
	return &StreamHandler{
		streamUseCase: streamUseCase,
		logger:        logger,
		errorHandler:  NewErrorHandler(logger),
	}
}

func (h *StreamHandler) GetStreamUseCase() usecase.StreamUseCaseInterface {
	return h.streamUseCase
}

func (h *StreamHandler) GetLogger() logger.Logger {
	return h.logger
}

func (h *StreamHandler) GetStreamData(c *fiber.Ctx) error {
	streamName := "todo-stream"
	countStr := c.Query("count", "10")
	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil || count <= 0 {
		return h.errorHandler.Handle(c, apperrors.NewValidationError("Invalid count parameter", err))
	}

	start := c.Query("start", "")

	messages, err := h.streamUseCase.GetStreamData(streamName, count, start)
	if err != nil {
		return h.errorHandler.Handle(c, err)
	}

	return dto.RespondWithJSON(c, fiber.StatusOK, messages)
}
