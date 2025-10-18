package redis_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ahmadrezamusthafa/ice-todo-service/adapter/persistence/redis"
	mock_redis_adapter "github.com/ahmadrezamusthafa/ice-todo-service/mock/adapter/persistence/redis"
	redisv8 "github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPublish(t *testing.T) {
	t.Log("Starting TestPublish")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStreamInterface := mock_redis_adapter.NewMockStreamInterface(ctrl)
	streamRepo := redis.NewStreamRepository(mockStreamInterface)

	testCases := []struct {
		name          string
		streamName    string
		data          map[string]interface{}
		mockBehavior  func()
		expectedID    string
		expectedError error
	}{
		{
			name:       "Success with string values",
			streamName: "test-stream",
			data: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			mockBehavior: func() {
				expectedValues := map[string]string{
					"key1": "value1",
					"key2": "value2",
				}
				cmd := redisv8.NewStringCmd(context.Background())
				cmd.SetVal("1234-0")
				mockStreamInterface.EXPECT().XAdd(gomock.Any(), &redisv8.XAddArgs{
					Stream: "test-stream",
					Values: expectedValues,
				}).Return(cmd)
			},
			expectedID:    "1234-0",
			expectedError: nil,
		},
		{
			name:       "Success with mixed value types",
			streamName: "test-stream",
			data: map[string]interface{}{
				"key1": "value1",
				"key2": 123,
				"key3": true,
			},
			mockBehavior: func() {
				expectedValues := map[string]string{
					"key1": "value1",
					"key2": "123",
					"key3": "true",
				}
				cmd := redisv8.NewStringCmd(context.Background())
				cmd.SetVal("1235-0")
				mockStreamInterface.EXPECT().XAdd(gomock.Any(), &redisv8.XAddArgs{
					Stream: "test-stream",
					Values: expectedValues,
				}).Return(cmd)
			},
			expectedID:    "1235-0",
			expectedError: nil,
		},
		{
			name:       "Redis Error",
			streamName: "test-stream",
			data: map[string]interface{}{
				"key1": "value1",
			},
			mockBehavior: func() {
				expectedValues := map[string]string{
					"key1": "value1",
				}
				cmd := redisv8.NewStringCmd(context.Background())
				cmd.SetErr(errors.New("redis connection error"))
				mockStreamInterface.EXPECT().XAdd(gomock.Any(), &redisv8.XAddArgs{
					Stream: "test-stream",
					Values: expectedValues,
				}).Return(cmd)
			},
			expectedID:    "",
			expectedError: errors.New("redis connection error"),
		},
		{
			name:       "Empty Data",
			streamName: "test-stream",
			data:       map[string]interface{}{},
			mockBehavior: func() {
				cmd := redisv8.NewStringCmd(context.Background())
				cmd.SetVal("1236-0")
				mockStreamInterface.EXPECT().XAdd(gomock.Any(), &redisv8.XAddArgs{
					Stream: "test-stream",
					Values: map[string]string{},
				}).Return(cmd)
			},
			expectedID:    "1236-0",
			expectedError: nil,
		},
	}

	t.Run("Publish", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {

				tc.mockBehavior()

				result, err := streamRepo.Publish(tc.streamName, tc.data)

				if tc.expectedError != nil {
					assert.Error(t, err)
					assert.Equal(t, tc.expectedError.Error(), err.Error())
					assert.Empty(t, result)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.expectedID, result)
				}
			})
		}
	})
}
