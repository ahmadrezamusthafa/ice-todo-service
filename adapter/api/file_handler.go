package api

import (
	"github.com/ahmadrezamusthafa/ice-todo-service/adapter/dto"
	"github.com/ahmadrezamusthafa/ice-todo-service/domain/validator"
	"github.com/ahmadrezamusthafa/ice-todo-service/infrastructure/logger"
	"github.com/ahmadrezamusthafa/ice-todo-service/usecase"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

type FileHandler struct {
	fileUseCase   *usecase.FileUseCase
	fileValidator *validator.FileValidator
	logger        logger.Logger
}

func NewFileHandler(fileUseCase *usecase.FileUseCase, maxFileSize int64, logger logger.Logger) *FileHandler {
	return &FileHandler{
		fileUseCase:   fileUseCase,
		fileValidator: validator.NewFileValidator(),
		logger:        logger,
	}
}

func (h *FileHandler) GetFileUseCase() *usecase.FileUseCase {
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
		h.logger.Error("Failed to get file from form: %v", err)
		return dto.RespondWithError(c, fiber.StatusBadRequest, "No file provided or invalid file")
	}

	if err := h.fileValidator.ValidateFileSize(file.Size); err != nil {
		return dto.RespondWithError(c, fiber.StatusBadRequest, err.Error())
	}

	ext := filepath.Ext(file.Filename)
	if err := h.fileValidator.ValidateFileType(ext); err != nil {
		return dto.RespondWithError(c, fiber.StatusBadRequest, err.Error())
	}

	src, err := file.Open()
	if err != nil {
		h.logger.Error("Failed to open uploaded file: %v", err)
		return dto.RespondWithError(c, fiber.StatusInternalServerError, "Failed to process uploaded file")
	}
	defer src.Close()

	fileID, err := h.fileUseCase.UploadFile(file.Filename, file.Size, src)
	if err != nil {
		h.logger.Error("Failed to upload file %s: %v", file.Filename, err)
		return dto.RespondWithError(c, fiber.StatusInternalServerError, "Failed to upload file: internal server error")
	}

	return dto.RespondWithJSON(c, fiber.StatusCreated, map[string]string{"fileId": fileID})
}

func (h *FileHandler) GetFile(c *fiber.Ctx) error {
	fileID := c.Params("id")
	if fileID == "" {
		return dto.RespondWithError(c, fiber.StatusBadRequest, "File ID is required")
	}

	fileContent, err := h.fileUseCase.GetFile(fileID)
	if err != nil {
		h.logger.Error("Failed to retrieve file with ID %s: %v", fileID, err)
		return dto.RespondWithError(c, fiber.StatusNotFound, "File not found or could not be retrieved")
	}

	c.Set("Content-Type", "application/octet-stream")

	return c.Send(fileContent)
}
