package mapper

import (
	"time"
)

type TodoResponse struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"dueDate"`
	FileID      string    `json:"fileId,omitempty"`
}

type TodoListResponse struct {
	Items []TodoResponse `json:"items"`
	Count int            `json:"count"`
}

type DeleteResponse struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}
