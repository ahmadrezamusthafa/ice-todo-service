package mapper

import (
	"github.com/ahmadrezamusthafa/ice-todo-service/adapter/dto"
	"github.com/ahmadrezamusthafa/ice-todo-service/domain/entity"
	"github.com/google/uuid"
)

func TodoRequestToEntity(req *dto.CreateTodoRequest) *entity.TodoItem {
	return entity.NewTodoItem(req.Description, req.DueDate, req.FileID)
}

func UpdateTodoRequestToEntity(req *dto.UpdateTodoRequest, existingTodo *entity.TodoItem) *entity.TodoItem {

	if req.Description != nil {
		existingTodo.Description = *req.Description
	}

	if req.DueDate != nil {
		existingTodo.DueDate = *req.DueDate
	}

	if req.FileID != "" {
		existingTodo.FileID = req.FileID
	}

	return existingTodo
}

func TodoEntityToResponse(todo *entity.TodoItem) *TodoResponse {
	return &TodoResponse{
		ID:          todo.ID.String(),
		Description: todo.Description,
		DueDate:     todo.DueDate,
		FileID:      todo.FileID,
	}
}

func TodoEntitiesToResponse(todos []*entity.TodoItem) *TodoListResponse {
	response := &TodoListResponse{
		Items: make([]TodoResponse, 0, len(todos)),
		Count: len(todos),
	}

	for _, todo := range todos {
		response.Items = append(response.Items, *TodoEntityToResponse(todo))
	}

	return response
}

func CreateDeleteResponse(id uuid.UUID, message string) *DeleteResponse {
	return &DeleteResponse{
		ID:      id.String(),
		Message: message,
	}
}
