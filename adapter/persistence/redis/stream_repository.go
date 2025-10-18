package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type StreamInterface interface {
	XAdd(ctx context.Context, a *redis.XAddArgs) *redis.StringCmd
	XRead(ctx context.Context, a *redis.XReadArgs) *redis.XStreamSliceCmd
}

type StreamRepository struct {
	client StreamInterface
}

func NewStreamRepository(client StreamInterface) *StreamRepository {
	return &StreamRepository{
		client: client,
	}
}

func (r *StreamRepository) Publish(streamName string, data map[string]interface{}) (string, error) {
	ctx := context.Background()

	values := make(map[string]string)
	for k, v := range data {
		switch val := v.(type) {
		case string:
			values[k] = val
		default:
			values[k] = fmt.Sprintf("%v", val)
		}
	}

	result := r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: values,
	})

	return result.Result()
}

func (r *StreamRepository) GetStream(streamName string, count int64, start string) ([]map[string]interface{}, error) {
	ctx := context.Background()
	if start == "" {
		start = "0"
	}

	streams := []string{streamName, start}
	result := r.client.XRead(ctx, &redis.XReadArgs{
		Streams: streams,
		Count:   count,
		Block:   0,
	})

	xStreams, err := result.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to read stream: %w", err)
	}

	messages := []map[string]interface{}{}
	for _, xStream := range xStreams {
		for _, xMessage := range xStream.Messages {
			msg := make(map[string]interface{})
			for k, v := range xMessage.Values {
				msg[k] = v
			}

			msg["id"] = xMessage.ID
			messages = append(messages, msg)
		}
	}

	return messages, nil
}
