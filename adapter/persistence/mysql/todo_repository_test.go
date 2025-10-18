package mysql_test

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ahmadrezamusthafa/ice-todo-service/adapter/persistence/mysql"
	"github.com/ahmadrezamusthafa/ice-todo-service/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := mysql.NewTodoRepository(db)

	testCases := []struct {
		name          string
		input         *entity.TodoItem
		mockBehavior  func()
		expectedError error
	}{
		{
			name: "Success",
			input: &entity.TodoItem{
				ID:          uuid.New(),
				Description: "Test Todo",
				DueDate:     time.Now().Add(24 * time.Hour),
				FileID:      "file123",
			},
			mockBehavior: func() {
				mock.ExpectExec("INSERT INTO todo_items").WithArgs(
					AnyUUID{},
					"Test Todo",
					sqlmock.AnyArg(),
					"file123",
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
		},
		{
			name: "Database Error",
			input: &entity.TodoItem{
				ID:          uuid.New(),
				Description: "Test Todo",
				DueDate:     time.Now().Add(24 * time.Hour),
				FileID:      "file123",
			},
			mockBehavior: func() {
				mock.ExpectExec("INSERT INTO todo_items").WithArgs(
					AnyUUID{},
					"Test Todo",
					sqlmock.AnyArg(),
					"file123",
				).WillReturnError(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	t.Run("Create", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {

				tc.mockBehavior()

				result, err := repo.Create(tc.input)

				if tc.expectedError != nil {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tc.expectedError.Error())
					assert.Nil(t, result)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.input, result)
				}

				assert.NoError(t, mock.ExpectationsWereMet())
			})
		}
	})
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := mysql.NewTodoRepository(db)

	todoID := uuid.New()
	testCases := []struct {
		name          string
		id            uuid.UUID
		mockBehavior  func()
		expectedTodo  *entity.TodoItem
		expectedError error
	}{
		{
			name: "Success",
			id:   todoID,
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "description", "due_date", "file_id"}).AddRow(
					todoID.String(),
					"Test Todo",
					time.Now().Add(24*time.Hour),
					"file123",
				)
				mock.ExpectQuery("SELECT id, description, due_date, file_id FROM todo_items WHERE id").WithArgs(
					todoID.String(),
				).WillReturnRows(rows)
			},
			expectedTodo: &entity.TodoItem{
				ID:          todoID,
				Description: "Test Todo",
				DueDate:     time.Now().Add(24 * time.Hour),
				FileID:      "file123",
			},
			expectedError: nil,
		},
		{
			name: "Not Found",
			id:   todoID,
			mockBehavior: func() {
				mock.ExpectQuery("SELECT id, description, due_date, file_id FROM todo_items WHERE id").WithArgs(
					todoID.String(),
				).WillReturnError(sql.ErrNoRows)
			},
			expectedTodo:  nil,
			expectedError: sql.ErrNoRows,
		},
		{
			name: "Database Error",
			id:   todoID,
			mockBehavior: func() {
				mock.ExpectQuery("SELECT id, description, due_date, file_id FROM todo_items WHERE id").WithArgs(
					todoID.String(),
				).WillReturnError(errors.New("database error"))
			},
			expectedTodo:  nil,
			expectedError: errors.New("database error"),
		},
	}

	t.Run("GetByID", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {

				tc.mockBehavior()

				result, err := repo.GetByID(tc.id)

				if tc.expectedError != nil {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tc.expectedError.Error())
					assert.Nil(t, result)
				} else {
					assert.NoError(t, err)

					assert.Equal(t, tc.expectedTodo.ID, result.ID)
					assert.Equal(t, tc.expectedTodo.Description, result.Description)
					assert.Equal(t, tc.expectedTodo.FileID, result.FileID)
				}

				assert.NoError(t, mock.ExpectationsWereMet())
			})
		}
	})
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := mysql.NewTodoRepository(db)

	todoID := uuid.New()
	testCases := []struct {
		name          string
		id            uuid.UUID
		input         *entity.TodoItem
		mockBehavior  func()
		expectedTodo  *entity.TodoItem
		expectedError error
	}{
		{
			name: "Success",
			id:   todoID,
			input: &entity.TodoItem{
				ID:          todoID,
				Description: "Updated Todo",
				DueDate:     time.Now().Add(48 * time.Hour),
				FileID:      "file456",
			},
			mockBehavior: func() {
				mock.ExpectExec("UPDATE todo_items SET description").WithArgs(
					"Updated Todo",
					sqlmock.AnyArg(),
					"file456",
					todoID.String(),
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedTodo: &entity.TodoItem{
				ID:          todoID,
				Description: "Updated Todo",
				DueDate:     time.Now().Add(48 * time.Hour),
				FileID:      "file456",
			},
			expectedError: nil,
		},
		{
			name: "Database Error",
			id:   todoID,
			input: &entity.TodoItem{
				ID:          todoID,
				Description: "Updated Todo",
				DueDate:     time.Now().Add(48 * time.Hour),
				FileID:      "file456",
			},
			mockBehavior: func() {
				mock.ExpectExec("UPDATE todo_items SET description").WithArgs(
					"Updated Todo",
					sqlmock.AnyArg(),
					"file456",
					todoID.String(),
				).WillReturnError(errors.New("database error"))
			},
			expectedTodo:  nil,
			expectedError: errors.New("database error"),
		},
	}

	t.Run("Update", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {

				tc.mockBehavior()

				result, err := repo.Update(tc.id, tc.input)

				if tc.expectedError != nil {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tc.expectedError.Error())
					assert.Nil(t, result)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.input, result)
				}

				assert.NoError(t, mock.ExpectationsWereMet())
			})
		}
	})
}

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := mysql.NewTodoRepository(db)

	todoID1 := uuid.New()
	todoID2 := uuid.New()
	testCases := []struct {
		name          string
		mockBehavior  func()
		expectedTodos []*entity.TodoItem
		expectedError error
	}{
		{
			name: "Success",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "description", "due_date", "file_id"}).
					AddRow(todoID1.String(), "Todo 1", time.Now().Add(24*time.Hour), "file123").
					AddRow(todoID2.String(), "Todo 2", time.Now().Add(48*time.Hour), "file456")

				mock.ExpectQuery("SELECT id, description, due_date, file_id FROM todo_items").WillReturnRows(rows)
			},
			expectedTodos: []*entity.TodoItem{
				{
					ID:          todoID1,
					Description: "Todo 1",
					DueDate:     time.Now().Add(24 * time.Hour),
					FileID:      "file123",
				},
				{
					ID:          todoID2,
					Description: "Todo 2",
					DueDate:     time.Now().Add(48 * time.Hour),
					FileID:      "file456",
				},
			},
			expectedError: nil,
		},
		{
			name: "Empty Result",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "description", "due_date", "file_id"})
				mock.ExpectQuery("SELECT id, description, due_date, file_id FROM todo_items").WillReturnRows(rows)
			},
			expectedTodos: []*entity.TodoItem{},
			expectedError: nil,
		},
		{
			name: "Database Error",
			mockBehavior: func() {
				mock.ExpectQuery("SELECT id, description, due_date, file_id FROM todo_items").WillReturnError(errors.New("database error"))
			},
			expectedTodos: nil,
			expectedError: errors.New("database error"),
		},
	}

	t.Run("GetAll", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {

				tc.mockBehavior()

				result, err := repo.GetAll()

				if tc.expectedError != nil {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tc.expectedError.Error())
					assert.Nil(t, result)
				} else {
					assert.NoError(t, err)
					if len(tc.expectedTodos) == 0 {
						assert.Empty(t, result)
					} else {
						assert.Equal(t, len(tc.expectedTodos), len(result))
						for i, todo := range result {
							assert.Equal(t, tc.expectedTodos[i].ID, todo.ID)
							assert.Equal(t, tc.expectedTodos[i].Description, todo.Description)
							assert.Equal(t, tc.expectedTodos[i].FileID, todo.FileID)
						}
					}
				}

				assert.NoError(t, mock.ExpectationsWereMet())
			})
		}
	})
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := mysql.NewTodoRepository(db)

	todoID := uuid.New()
	testCases := []struct {
		name          string
		id            uuid.UUID
		mockBehavior  func()
		expectedError error
	}{
		{
			name: "Success",
			id:   todoID,
			mockBehavior: func() {
				mock.ExpectExec("DELETE FROM todo_items WHERE id").WithArgs(
					todoID.String(),
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
		},
		{
			name: "Not Found",
			id:   todoID,
			mockBehavior: func() {
				mock.ExpectExec("DELETE FROM todo_items WHERE id").WithArgs(
					todoID.String(),
				).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: errors.New("not found"),
		},
		{
			name: "Database Error",
			id:   todoID,
			mockBehavior: func() {
				mock.ExpectExec("DELETE FROM todo_items WHERE id").WithArgs(
					todoID.String(),
				).WillReturnError(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	t.Run("Delete", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {

				tc.mockBehavior()

				err := repo.Delete(tc.id)

				if tc.expectedError != nil {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tc.expectedError.Error())
				} else {
					assert.NoError(t, err)
				}

				assert.NoError(t, mock.ExpectationsWereMet())
			})
		}
	})
}

type AnyUUID struct{}

func (a AnyUUID) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}
