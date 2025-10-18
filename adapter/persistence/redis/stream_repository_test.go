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
					assert.Contains(t, err.Error(), tc.expectedError.Error())
					assert.Empty(t, result)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.expectedID, result)
				}
			})
		}
	})
}

func TestGetStream(t *testing.T) {
	t.Log("Starting TestGetStream")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStreamInterface := mock_redis_adapter.NewMockStreamInterface(ctrl)
	streamRepo := redis.NewStreamRepository(mockStreamInterface)

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
			name:       "Success with multiple messages",
			streamName: "test-stream",
			count:      10,
			start:      "0",
			mockBehavior: func() {

				cmd := redisv8.NewXStreamSliceCmd(context.Background())

				messages := []redisv8.XMessage{
					{
						ID: "1234-0",
						Values: map[string]interface{}{
							"key1": "value1",
							"key2": "value2",
						},
					},
					{
						ID: "1235-0",
						Values: map[string]interface{}{
							"key3": "value3",
						},
					},
				}

				streams := []redisv8.XStream{
					{
						Stream:   "test-stream",
						Messages: messages,
					},
				}

				cmd.SetVal(streams)

				mockStreamInterface.EXPECT().XRead(gomock.Any(), &redisv8.XReadArgs{
					Streams: []string{"test-stream", "0"},
					Count:   10,
					Block:   0,
				}).Return(cmd)
			},
			expectedData: []map[string]interface{}{
				{
					"key1": "value1",
					"key2": "value2",
					"id":   "1234-0",
				},
				{
					"key3": "value3",
					"id":   "1235-0",
				},
			},
			expectedError: nil,
		},
		{
			name:       "Empty start parameter",
			streamName: "test-stream",
			count:      5,
			start:      "",
			mockBehavior: func() {
				cmd := redisv8.NewXStreamSliceCmd(context.Background())
				cmd.SetVal([]redisv8.XStream{})

				mockStreamInterface.EXPECT().XRead(gomock.Any(), &redisv8.XReadArgs{
					Streams: []string{"test-stream", "0"},
					Count:   5,
					Block:   0,
				}).Return(cmd)
			},
			expectedData:  []map[string]interface{}{},
			expectedError: nil,
		},
		{
			name:       "Redis error",
			streamName: "test-stream",
			count:      10,
			start:      "0",
			mockBehavior: func() {
				cmd := redisv8.NewXStreamSliceCmd(context.Background())
				cmd.SetErr(errors.New("redis connection error"))

				mockStreamInterface.EXPECT().XRead(gomock.Any(), gomock.Any()).Return(cmd)
			},
			expectedData:  nil,
			expectedError: errors.New("failed to read stream: redis connection error"),
		},
	}

	t.Run("GetStream", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {

				tc.mockBehavior()

				result, err := streamRepo.GetStream(tc.streamName, tc.count, tc.start)

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
	})

}
