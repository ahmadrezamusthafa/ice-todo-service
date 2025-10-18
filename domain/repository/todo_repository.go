package repository

import (
	"github.com/ahmadrezamusthafa/ice-todo-service/domain/entity"
	"github.com/google/uuid"
)

type TodoRepository interface {
	Create(todo *entity.TodoItem) (*entity.TodoItem, error)
	GetByID(id uuid.UUID) (*entity.TodoItem, error)
	Update(id uuid.UUID, todo *entity.TodoItem) (*entity.TodoItem, error)
	GetAll() ([]*entity.TodoItem, error)
	Delete(id uuid.UUID) error
}
