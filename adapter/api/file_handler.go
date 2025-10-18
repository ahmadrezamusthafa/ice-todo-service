package api

import (
	"github.com/ahmadrezamusthafa/ice-todo-service/adapter/dto"
	"github.com/ahmadrezamusthafa/ice-todo-service/adapter/dto/mapper"
	"github.com/ahmadrezamusthafa/ice-todo-service/domain/apperrors"
	"github.com/ahmadrezamusthafa/ice-todo-service/domain/validator"
	"github.com/ahmadrezamusthafa/ice-todo-service/infrastructure/logger"
	"github.com/ahmadrezamusthafa/ice-todo-service/usecase"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

type FileHandler struct {
	fileUseCase   usecase.FileUseCaseInterface
	fileValidator *validator.FileValidator
	logger        logger.Logger
	errorHandler  *ErrorHandler
}

func NewFileHandler(fileUseCase usecase.FileUseCaseInterface, maxFileSize int64, logger logger.Logger) *FileHandler {
	return &FileHandler{
		fileUseCase:   fileUseCase,
		fileValidator: validator.NewFileValidator(),
		logger:        logger,
		errorHandler:  NewErrorHandler(logger),
	}
}

func (h *FileHandler) GetFileUseCase() usecase.FileUseCaseInterface {
	return h.fileUseCase
}

func (h *FileHandler) GetFileValidator() *validator.FileValidator {
	return h.fileValidator
}

func (h *FileHandler) GetLogger() logger.Logger {
	return h.logger
}

func (h *FileHandler) UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return h.errorHandler.Handle(c, apperrors.NewValidationError("No file provided or invalid file", err))
	}

	if err := h.fileValidator.ValidateFileSize(file.Size); err != nil {
		return h.errorHandler.Handle(c, apperrors.NewValidationError("Invalid file size", err))
	}

	ext := filepath.Ext(file.Filename)
	if err := h.fileValidator.ValidateFileType(ext); err != nil {
		return h.errorHandler.Handle(c, apperrors.NewValidationError("Invalid file type", err))
	}

	src, err := file.Open()
	if err != nil {
		return h.errorHandler.Handle(c, apperrors.NewInternalError("Failed to process uploaded file", err))
	}
	defer src.Close()

	fileID, err := h.fileUseCase.UploadFile(file.Filename, file.Size, src)
	if err != nil {
		return h.errorHandler.Handle(c, err)
	}

	response := mapper.FileToUploadResponse(fileID, file)
	return dto.RespondWithJSON(c, fiber.StatusCreated, response)
}

func (h *FileHandler) GetFile(c *fiber.Ctx) error {
	fileID := c.Params("id")
	if fileID == "" {
		return h.errorHandler.Handle(c, apperrors.NewValidationError("File ID is required", nil))
	}

	fileContent, err := h.fileUseCase.GetFile(fileID)
	if err != nil {
		return h.errorHandler.Handle(c, err)
	}

	c.Set("Content-Type", "application/octet-stream")
	return c.Send(fileContent)
}
