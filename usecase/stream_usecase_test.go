package usecase_test

import (
	"errors"
	mock_infrastructure "github.com/ahmadrezamusthafa/ice-todo-service/mock/infrastructure"
	"testing"

	"github.com/ahmadrezamusthafa/ice-todo-service/mock/repository"
	"github.com/ahmadrezamusthafa/ice-todo-service/usecase"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetStreamData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStreamRepo := mock_repository.NewMockStreamRepository(ctrl)
	mockLogger := mock_infrastructure.NewMockLogger(ctrl)

	streamUseCase := usecase.NewStreamUseCase(mockStreamRepo, mockLogger)

	testCases := []struct {
		name          string
		streamName    string
		count         int64
		start         string
		mockBehavior  func()
		expectedData  []map[string]interface{}
		expectedError error
	}{
		{
			name:       "Success",
			streamName: "test-stream",
			count:      10,
			start:      "0",
			mockBehavior: func() {
				expectedData := []map[string]interface{}{
					{
						"key1": "value1",
						"id":   "1234-0",
					},
					{
						"key2": "value2",
						"id":   "1235-0",
					},
				}

				mockLogger.EXPECT().Info("Getting data from stream: %s (count: %d, start: %s)", "test-stream", int64(10), "0").Times(1)
				mockLogger.EXPECT().Info("Retrieved %d messages from stream", 2).Times(1)

				mockStreamRepo.EXPECT().GetStream(
					"test-stream",
					int64(10),
					"0",
				).Return(expectedData, nil)
			},
			expectedData: []map[string]interface{}{
				{
					"key1": "value1",
					"id":   "1234-0",
				},
				{
					"key2": "value2",
					"id":   "1235-0",
				},
			},
			expectedError: nil,
		},
		{
			name:       "Empty Result",
			streamName: "test-stream",
			count:      5,
			start:      "",
			mockBehavior: func() {

				mockLogger.EXPECT().Info("Getting data from stream: %s (count: %d, start: %s)", "test-stream", int64(5), "").Times(1)
				mockLogger.EXPECT().Info("Retrieved %d messages from stream", 0).Times(1)

				mockStreamRepo.EXPECT().GetStream(
					"test-stream",
					int64(5),
					"",
				).Return([]map[string]interface{}{}, nil)
			},
			expectedData:  []map[string]interface{}{},
			expectedError: nil,
		},
		{
			name:       "Repository Error",
			streamName: "test-stream",
			count:      10,
			start:      "0",
			mockBehavior: func() {

				mockLogger.EXPECT().Info("Getting data from stream: %s (count: %d, start: %s)", "test-stream", int64(10), "0").Times(1)
				mockLogger.EXPECT().Error("Failed to get stream data: %v", gomock.Any()).Times(1)

				mockStreamRepo.EXPECT().GetStream(
					"test-stream",
					int64(10),
					"0",
				).Return(nil, errors.New("repository error"))
			},
			expectedData:  nil,
			expectedError: errors.New("repository error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior()

			result, err := streamUseCase.GetStreamData(tc.streamName, tc.count, tc.start)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedData, result)
			}
		})
	}
}
