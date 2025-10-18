package usecase_test

import (
	"errors"
	"github.com/ahmadrezamusthafa/ice-todo-service/domain/apperrors"
	mock_infrastructure "github.com/ahmadrezamusthafa/ice-todo-service/mock/infrastructure"
	mock_repository "github.com/ahmadrezamusthafa/ice-todo-service/mock/repository"
	"testing"
	"time"

	"github.com/ahmadrezamusthafa/ice-todo-service/domain/entity"
	"github.com/ahmadrezamusthafa/ice-todo-service/usecase"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateTodo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTodoRepo := mock_repository.NewMockTodoRepository(ctrl)
	mockStreamRepo := mock_repository.NewMockStreamRepository(ctrl)
	mockLogger := mock_infrastructure.NewMockLogger(ctrl)

	todoUseCase := usecase.NewTodoUseCase(mockTodoRepo, mockStreamRepo, mockLogger)

	testCases := []struct {
		name          string
		input         *entity.TodoItem
		expectedTodo  *entity.TodoItem
		repoError     error
		streamError   error
		expectedError error
	}{
		{
			name: "Success",
			input: &entity.TodoItem{
				ID:          uuid.New(),
				Description: "Test Todo",
				DueDate:     time.Now().Add(24 * time.Hour),
			},
			expectedTodo: &entity.TodoItem{
				ID:          uuid.New(),
				Description: "Test Todo",
				DueDate:     time.Now().Add(24 * time.Hour),
			},
			repoError:     nil,
			streamError:   nil,
			expectedError: nil,
		},
		{
			name: "Repository Error",
			input: &entity.TodoItem{
				ID:          uuid.New(),
				Description: "Test Todo",
				DueDate:     time.Now().Add(24 * time.Hour),
			},
			expectedTodo:  nil,
			repoError:     errors.New("repository error"),
			streamError:   nil,
			expectedError: errors.New("repository error"),
		},
		{
			name: "Stream Error But Todo Created",
			input: &entity.TodoItem{
				ID:          uuid.New(),
				Description: "Test Todo",
				DueDate:     time.Now().Add(24 * time.Hour),
			},
			expectedTodo: &entity.TodoItem{
				ID:          uuid.New(),
				Description: "Test Todo",
				DueDate:     time.Now().Add(24 * time.Hour),
			},
			repoError:     nil,
			streamError:   errors.New("stream error"),
			expectedError: nil,
		},
	}

	t.Run("CreateTodo", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {

				mockTodoRepo.EXPECT().
					Create(gomock.Eq(tc.input)).
					Return(tc.expectedTodo, tc.repoError).AnyTimes()

				if tc.repoError == nil {

					expectedData := map[string]interface{}{
						"id":          tc.expectedTodo.ID.String(),
						"description": tc.expectedTodo.Description,
						"dueDate":     tc.expectedTodo.DueDate.Format("2006-01-02T15:04:05Z07:00"),
						"fileId":      tc.expectedTodo.FileID,
					}

					mockStreamRepo.EXPECT().
						Publish("todo-stream", gomock.Eq(expectedData)).
						Return("", tc.streamError).AnyTimes()

					if tc.streamError != nil {
						mockLogger.EXPECT().
							Error(gomock.Any(), gomock.Any()).AnyTimes()
					}
				}

				result, err := todoUseCase.CreateTodo(tc.input)

				if tc.expectedError != nil {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tc.expectedError.Error())
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.expectedTodo, result)
				}
			})
		}
	})
}

func TestUpdateTodo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTodoRepo := mock_repository.NewMockTodoRepository(ctrl)
	mockStreamRepo := mock_repository.NewMockStreamRepository(ctrl)
	mockLogger := mock_infrastructure.NewMockLogger(ctrl)

	todoUseCase := usecase.NewTodoUseCase(mockTodoRepo, mockStreamRepo, mockLogger)

	testCases := []struct {
		name          string
		id            uuid.UUID
		input         *entity.TodoItem
		expectedTodo  *entity.TodoItem
		repoError     error
		streamError   error
		expectedError error
	}{
		{
			name: "Success",
			id:   uuid.New(),
			input: &entity.TodoItem{
				Description: "Updated Todo",
				DueDate:     time.Now().Add(48 * time.Hour),
			},
			expectedTodo: &entity.TodoItem{
				ID:          uuid.New(),
				Description: "Updated Todo",
				DueDate:     time.Now().Add(48 * time.Hour),
			},
			repoError:     nil,
			streamError:   nil,
			expectedError: nil,
		},
		{
			name: "Repository Error",
			id:   uuid.New(),
			input: &entity.TodoItem{
				Description: "Updated Todo",
				DueDate:     time.Now().Add(48 * time.Hour),
			},
			expectedTodo:  nil,
			repoError:     errors.New("repository error"),
			streamError:   nil,
			expectedError: errors.New("repository error"),
		},
		{
			name: "Stream Error But Todo Updated",
			id:   uuid.New(),
			input: &entity.TodoItem{
				Description: "Updated Todo",
				DueDate:     time.Now().Add(48 * time.Hour),
			},
			expectedTodo: &entity.TodoItem{
				ID:          uuid.New(),
				Description: "Updated Todo",
				DueDate:     time.Now().Add(48 * time.Hour),
			},
			repoError:     nil,
			streamError:   errors.New("stream error"),
			expectedError: nil,
		},
	}

	t.Run("UpdateTodo", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {

				mockTodoRepo.EXPECT().
					Update(gomock.Eq(tc.id), gomock.Eq(tc.input)).
					Return(tc.expectedTodo, tc.repoError).AnyTimes()

				if tc.repoError == nil {
					expectedData := map[string]interface{}{
						"id":          tc.expectedTodo.ID.String(),
						"description": tc.expectedTodo.Description,
						"dueDate":     tc.expectedTodo.DueDate.Format("2006-01-02T15:04:05Z07:00"),
						"fileId":      tc.expectedTodo.FileID,
					}

					mockStreamRepo.EXPECT().
						Publish("todo-stream", gomock.Eq(expectedData)).
						Return("", tc.streamError).AnyTimes()

					if tc.streamError != nil {
						mockLogger.EXPECT().
							Error(gomock.Any(), gomock.Any()).AnyTimes()
					}
				}

				result, err := todoUseCase.UpdateTodo(tc.id, tc.input)

				if tc.expectedError != nil {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tc.expectedError.Error())
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.expectedTodo, result)
				}
			})
		}
	})
}

func TestGetTodoByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTodoRepo := mock_repository.NewMockTodoRepository(ctrl)
	mockStreamRepo := mock_repository.NewMockStreamRepository(ctrl)
	mockLogger := mock_infrastructure.NewMockLogger(ctrl)

	todoUseCase := usecase.NewTodoUseCase(mockTodoRepo, mockStreamRepo, mockLogger)

	testCases := []struct {
		name          string
		id            uuid.UUID
		expectedTodo  *entity.TodoItem
		repoError     error
		expectedError error
	}{
		{
			name: "Success",
			id:   uuid.New(),
			expectedTodo: &entity.TodoItem{
				ID:          uuid.New(),
				Description: "Test Todo",
				DueDate:     time.Now(),
			},
			repoError:     nil,
			expectedError: nil,
		},
		{
			name:          "Not Found",
			id:            uuid.New(),
			expectedTodo:  nil,
			repoError:     errors.New("todo not found"),
			expectedError: errors.New("todo not found"),
		},
	}

	t.Run("GetTodoByID", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {

				mockTodoRepo.EXPECT().
					GetByID(gomock.Eq(tc.id)).
					Return(tc.expectedTodo, tc.repoError).AnyTimes()

				result, err := todoUseCase.GetTodoByID(tc.id)

				if tc.expectedError != nil {
					assert.Error(t, err)
					assert.Equal(t, tc.expectedError.Error(), err.Error())
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.expectedTodo, result)
				}
			})
		}
	})
}

func TestGetAllTodos(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTodoRepo := mock_repository.NewMockTodoRepository(ctrl)
	mockStreamRepo := mock_repository.NewMockStreamRepository(ctrl)
	mockLogger := mock_infrastructure.NewMockLogger(ctrl)

	todoUseCase := usecase.NewTodoUseCase(mockTodoRepo, mockStreamRepo, mockLogger)

	testCases := []struct {
		name          string
		expectedTodos []*entity.TodoItem
		repoError     error
		expectedError error
	}{
		{
			name: "Success",
			expectedTodos: []*entity.TodoItem{
				{
					ID:          uuid.New(),
					Description: "Todo 1",
					DueDate:     time.Now(),
				},
				{
					ID:          uuid.New(),
					Description: "Todo 2",
					DueDate:     time.Now().Add(24 * time.Hour),
				},
			},
			repoError:     nil,
			expectedError: nil,
		},
		{
			name:          "Repository Error",
			expectedTodos: nil,
			repoError:     errors.New("database error"),
			expectedError: errors.New("database error"),
		},
		{
			name:          "Empty List",
			expectedTodos: []*entity.TodoItem{},
			repoError:     nil,
			expectedError: nil,
		},
	}

	t.Run("GetAllTodos", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {

				mockTodoRepo.EXPECT().
					GetAll().
					Return(tc.expectedTodos, tc.repoError)

				result, err := todoUseCase.GetAllTodos()

				if tc.expectedError != nil {
					assert.Error(t, err)
					assert.Equal(t, tc.expectedError.Error(), err.Error())
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.expectedTodos, result)
				}
			})
		}
	})
}

func TestDeleteTodo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTodoRepo := mock_repository.NewMockTodoRepository(ctrl)
	mockStreamRepo := mock_repository.NewMockStreamRepository(ctrl)
	mockLogger := mock_infrastructure.NewMockLogger(ctrl)

	todoUseCase := usecase.NewTodoUseCase(mockTodoRepo, mockStreamRepo, mockLogger)

	testCases := []struct {
		name          string
		id            uuid.UUID
		todo          *entity.TodoItem
		getError      error
		deleteError   error
		streamError   error
		expectedError error
	}{
		{
			name: "Success",
			id:   uuid.New(),
			todo: &entity.TodoItem{
				ID:          uuid.New(),
				Description: "Todo to delete",
				DueDate:     time.Now(),
			},
			getError:      nil,
			deleteError:   nil,
			streamError:   nil,
			expectedError: nil,
		},
		{
			name:          "Todo Not Found",
			id:            uuid.New(),
			todo:          nil,
			getError:      errors.New("todo not found"),
			deleteError:   nil,
			streamError:   nil,
			expectedError: errors.New("todo not found"),
		},
		{
			name: "Delete Error",
			id:   uuid.New(),
			todo: &entity.TodoItem{
				ID:          uuid.New(),
				Description: "Todo to delete",
				DueDate:     time.Now(),
			},
			getError:      nil,
			deleteError:   errors.New("delete error"),
			streamError:   nil,
			expectedError: apperrors.NewInternalError("Failed to delete todo", errors.New("delete error")),
		},
		{
			name: "Stream Error But Todo Deleted",
			id:   uuid.New(),
			todo: &entity.TodoItem{
				ID:          uuid.New(),
				Description: "Todo to delete",
				DueDate:     time.Now(),
			},
			getError:      nil,
			deleteError:   nil,
			streamError:   errors.New("stream error"),
			expectedError: nil,
		},
	}

	t.Run("DeleteTodo", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {

				mockTodoRepo.EXPECT().
					GetByID(gomock.Eq(tc.id)).
					Return(tc.todo, tc.getError).AnyTimes()

				if tc.getError == nil {
					mockTodoRepo.EXPECT().
						Delete(gomock.Eq(tc.id)).
						Return(tc.deleteError).AnyTimes()

					if tc.deleteError == nil {
						expectedData := map[string]interface{}{
							"id":          tc.todo.ID.String(),
							"description": tc.todo.Description,
							"dueDate":     tc.todo.DueDate.Format("2006-01-02T15:04:05Z07:00"),
							"fileId":      tc.todo.FileID,
							"deleted":     true,
						}

						mockStreamRepo.EXPECT().
							Publish("todo-stream", gomock.Eq(expectedData)).
							Return("", tc.streamError).AnyTimes()

						if tc.streamError != nil {
							mockLogger.EXPECT().
								Error(gomock.Any(), gomock.Any()).AnyTimes()
						}
					}
				}

				err := todoUseCase.DeleteTodo(tc.id)

				if tc.expectedError != nil {
					assert.Error(t, err)
					assert.Equal(t, tc.expectedError.Error(), err.Error())
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}
