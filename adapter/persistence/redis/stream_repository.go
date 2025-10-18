package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type StreamInterface interface {
	XAdd(ctx context.Context, a *redis.XAddArgs) *redis.StringCmd
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
