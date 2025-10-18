package usecase_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	mock_infrastructure "github.com/ahmadrezamusthafa/ice-todo-service/mock/infrastructure"
	mock_repository "github.com/ahmadrezamusthafa/ice-todo-service/mock/repository"
	"github.com/ahmadrezamusthafa/ice-todo-service/usecase"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUploadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFileRepo := mock_repository.NewMockFileRepository(ctrl)
	mockLogger := mock_infrastructure.NewMockLogger(ctrl)

	fileUseCase := usecase.NewFileUseCase(mockFileRepo, mockLogger)

	testCases := []struct {
		name           string
		fileName       string
		fileSize       int64
		fileContent    io.Reader
		expectedFileID string
		repoError      error
		expectedError  error
	}{
		{
			name:           "Success",
			fileName:       "test.txt",
			fileSize:       100,
			fileContent:    bytes.NewReader([]byte("test content")),
			expectedFileID: "file123",
			repoError:      nil,
			expectedError:  nil,
		},
		{
			name:           "Repository Error",
			fileName:       "error.txt",
			fileSize:       50,
			fileContent:    bytes.NewReader([]byte("error content")),
			expectedFileID: "",
			repoError:      errors.New("storage error"),
			expectedError:  errors.New("failed to upload file: storage error"),
		},
	}

	t.Run("UploadFile", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {

				mockLogger.EXPECT().Info("Uploading file: %s (size: %d bytes)", tc.fileName, tc.fileSize)

				mockFileRepo.EXPECT().
					Upload(tc.fileName, tc.fileSize, gomock.Any()).
					Return(tc.expectedFileID, tc.repoError)

				if tc.repoError != nil {
					mockLogger.EXPECT().Error("Failed to upload file: %v", tc.repoError)
				} else {
					mockLogger.EXPECT().Info("File uploaded successfully with ID: %s", tc.expectedFileID)
				}

				fileID, err := fileUseCase.UploadFile(tc.fileName, tc.fileSize, tc.fileContent)

				if tc.expectedError != nil {
					assert.Error(t, err)
					assert.Equal(t, tc.expectedError.Error(), err.Error())
				} else {
					assert.NoError(t, err)
				}
				assert.Equal(t, tc.expectedFileID, fileID)
			})
		}
	})
}

func TestGetFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFileRepo := mock_repository.NewMockFileRepository(ctrl)
	mockLogger := mock_infrastructure.NewMockLogger(ctrl)

	fileUseCase := usecase.NewFileUseCase(mockFileRepo, mockLogger)

	testCases := []struct {
		name            string
		fileID          string
		expectedContent []byte
		repoError       error
		expectedError   error
	}{
		{
			name:            "Success",
			fileID:          "file123",
			expectedContent: []byte("file content"),
			repoError:       nil,
			expectedError:   nil,
		},
		{
			name:            "Repository Error",
			fileID:          "invalid-id",
			expectedContent: nil,
			repoError:       errors.New("file not found"),
			expectedError:   errors.New("failed to retrieve file: file not found"),
		},
	}

	t.Run("GetFile", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {

				mockLogger.EXPECT().Info("Retrieving file with ID: %s", tc.fileID)

				mockFileRepo.EXPECT().
					Get(tc.fileID).
					Return(tc.expectedContent, tc.repoError)

				if tc.repoError != nil {
					mockLogger.EXPECT().Error("Failed to retrieve file: %v", tc.repoError)
				} else {
					mockLogger.EXPECT().Info("File retrieved successfully, size: %d bytes", len(tc.expectedContent))
				}

				content, err := fileUseCase.GetFile(tc.fileID)

				if tc.expectedError != nil {
					assert.Error(t, err)
					assert.Equal(t, tc.expectedError.Error(), err.Error())
				} else {
					assert.NoError(t, err)
				}
				assert.Equal(t, tc.expectedContent, content)
			})
		}
	})
}
