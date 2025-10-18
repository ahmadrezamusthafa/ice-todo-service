package usecase

import (
	"github.com/ahmadrezamusthafa/ice-todo-service/domain/repository"
	"github.com/ahmadrezamusthafa/ice-todo-service/infrastructure/logger"
)

type StreamUseCaseInterface interface {
	GetStreamData(streamName string, count int64, start string) ([]map[string]interface{}, error)
}

type StreamUseCase struct {
	streamRepo repository.StreamRepository
	logger     logger.Logger
}

func NewStreamUseCase(streamRepo repository.StreamRepository, logger logger.Logger) StreamUseCaseInterface {
	return &StreamUseCase{
		streamRepo: streamRepo,
		logger:     logger,
	}
}

func (uc *StreamUseCase) GetStreamData(streamName string, count int64, start string) ([]map[string]interface{}, error) {
	uc.logger.Info("Getting data from stream: %s (count: %d, start: %s)", streamName, count, start)

	messages, err := uc.streamRepo.GetStream(streamName, count, start)
	if err != nil {
		uc.logger.Error("Failed to get stream data: %v", err)
		return nil, err
	}

	uc.logger.Info("Retrieved %d messages from stream", len(messages))
	return messages, nil
}
