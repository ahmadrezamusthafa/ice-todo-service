package dto

import (
	"time"
)

type CreateTodoRequest struct {
	Description string    `json:"description"`
	DueDate     time.Time `json:"dueDate"`
	FileID      string    `json:"fileId,omitempty"`
}

type UpdateTodoRequest struct {
	Description *string    `json:"description,omitempty"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
}
