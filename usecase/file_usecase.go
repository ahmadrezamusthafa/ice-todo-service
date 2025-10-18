package usecase

import (
	"github.com/ahmadrezamusthafa/ice-todo-service/domain/repository"
	"github.com/ahmadrezamusthafa/ice-todo-service/infrastructure/logger"
	"io"
)

type FileUseCaseInterface interface {
	UploadFile(fileName string, fileSize int64, fileContent io.Reader) (string, error)
	GetFile(fileID string) ([]byte, error)
}

type FileUseCase struct {
	fileRepo repository.FileRepository
	logger   logger.Logger
}

func NewFileUseCase(fileRepo repository.FileRepository, logger logger.Logger) FileUseCaseInterface {
	return &FileUseCase{
		fileRepo: fileRepo,
		logger:   logger,
	}
}

func (uc *FileUseCase) UploadFile(fileName string, fileSize int64, fileContent io.Reader) (string, error) {
	uc.logger.Info("Uploading file: %s (size: %d bytes)", fileName, fileSize)

	fileID, err := uc.fileRepo.Upload(fileName, fileSize, fileContent)
	if err != nil {
		uc.logger.Error("Failed to upload file: %v", err)
		return "", err
	}

	uc.logger.Info("File uploaded successfully with ID: %s", fileID)
	return fileID, nil
}

func (uc *FileUseCase) GetFile(fileID string) ([]byte, error) {
	uc.logger.Info("Retrieving file with ID: %s", fileID)

	fileContent, err := uc.fileRepo.Get(fileID)
	if err != nil {
		uc.logger.Error("Failed to retrieve file: %v", err)
		return nil, err
	}

	uc.logger.Info("File retrieved successfully, size: %d bytes", len(fileContent))
	return fileContent, nil
}
