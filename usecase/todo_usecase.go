package usecase

import (
	"github.com/ahmadrezamusthafa/ice-todo-service/domain/entity"
	"github.com/ahmadrezamusthafa/ice-todo-service/domain/repository"
	"github.com/ahmadrezamusthafa/ice-todo-service/infrastructure/logger"
	"github.com/google/uuid"
)

type TodoUseCase struct {
	todoRepo   repository.TodoRepository
	streamRepo repository.StreamRepository
	logger     logger.Logger
}

func NewTodoUseCase(todoRepo repository.TodoRepository, streamRepo repository.StreamRepository, logger logger.Logger) *TodoUseCase {
	return &TodoUseCase{
		todoRepo:   todoRepo,
		streamRepo: streamRepo,
		logger:     logger,
	}
}

func (uc *TodoUseCase) CreateTodo(todo *entity.TodoItem) (*entity.TodoItem, error) {
	createdTodo, err := uc.todoRepo.Create(todo)
	if err != nil {
		return nil, err
	}

	_, err = uc.publishToStream(createdTodo)
	if err != nil {
		uc.logger.Error("Failed to publish todo to stream: %v", err)
	}

	return createdTodo, nil
}

func (uc *TodoUseCase) UpdateTodo(id uuid.UUID, todo *entity.TodoItem) (*entity.TodoItem, error) {
	updatedTodo, err := uc.todoRepo.Update(id, todo)
	if err != nil {
		return nil, err
	}

	_, err = uc.publishToStream(updatedTodo)
	if err != nil {
		uc.logger.Error("Failed to publish updated todo to stream: %v", err)
	}

	return updatedTodo, nil
}

func (uc *TodoUseCase) GetTodoByID(id uuid.UUID) (*entity.TodoItem, error) {
	return uc.todoRepo.GetByID(id)
}

func (uc *TodoUseCase) publishToStream(todo *entity.TodoItem) (string, error) {
	data := map[string]interface{}{
		"id":          todo.ID.String(),
		"description": todo.Description,
		"dueDate":     todo.DueDate.Format("2006-01-02T15:04:05Z07:00"),
		"fileId":      todo.FileID,
	}

	return uc.streamRepo.Publish("todo-stream", data)
}
