package mysql

import (
	"database/sql"
	"github.com/ahmadrezamusthafa/ice-todo-service/domain/entity"
	"time"

	"github.com/google/uuid"
)

type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{
		db: db,
	}
}

func (r *TodoRepository) Create(todo *entity.TodoItem) (*entity.TodoItem, error) {
	query := `INSERT INTO todo_items (id, description, due_date, file_id) VALUES (?, ?, ?, ?)`

	_, err := r.db.Exec(query, todo.ID.String(), todo.Description, todo.DueDate, todo.FileID)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (r *TodoRepository) GetByID(id uuid.UUID) (*entity.TodoItem, error) {
	query := `SELECT id, description, due_date, file_id FROM todo_items WHERE id = ?`

	row := r.db.QueryRow(query, id.String())

	var idStr string
	var description string
	var dueDate time.Time
	var fileID string

	err := row.Scan(&idStr, &description, &dueDate, &fileID)
	if err != nil {
		return nil, err
	}

	todoID, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	return &entity.TodoItem{
		ID:          todoID,
		Description: description,
		DueDate:     dueDate,
		FileID:      fileID,
	}, nil
}

func (r *TodoRepository) Update(id uuid.UUID, todo *entity.TodoItem) (*entity.TodoItem, error) {
	query := `UPDATE todo_items SET description = ?, due_date = ? WHERE id = ?`

	_, err := r.db.Exec(query, todo.Description, todo.DueDate, id.String())
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (r *TodoRepository) GetAll() ([]*entity.TodoItem, error) {
	query := `SELECT id, description, due_date, file_id FROM todo_items`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*entity.TodoItem

	for rows.Next() {
		var idStr string
		var description string
		var dueDate time.Time
		var fileID string

		err := rows.Scan(&idStr, &description, &dueDate, &fileID)
		if err != nil {
			return nil, err
		}

		todoID, err := uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}

		todos = append(todos, &entity.TodoItem{
			ID:          todoID,
			Description: description,
			DueDate:     dueDate,
			FileID:      fileID,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *TodoRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM todo_items WHERE id = ?`

	result, err := r.db.Exec(query, id.String())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
