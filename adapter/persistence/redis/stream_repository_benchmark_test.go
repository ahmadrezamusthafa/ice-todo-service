package redis_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ahmadrezamusthafa/ice-todo-service/adapter/persistence/redis"
	"github.com/ahmadrezamusthafa/ice-todo-service/config"
	redisv8 "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

func setupTestRedis() (*redisv8.Client, error) {
	cfg := config.NewConfig()

	client := redisv8.NewClient(&redisv8.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       0,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return client, nil
}

func cleanupTestData(client *redisv8.Client, streamName string) {
	ctx := context.Background()
	_ = client.Del(ctx, streamName).Err()
}

func BenchmarkStreamRepositoryPublish(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}

	client, err := setupTestRedis()
	if err != nil {
		b.Fatalf("Failed to setup test Redis: %v", err)
	}
	defer client.Close()

	repo := redis.NewStreamRepository(client)

	streamName := fmt.Sprintf("benchmark-stream-%s", uuid.New().String())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		data := map[string]interface{}{
			"id":          uuid.New().String(),
			"description": fmt.Sprintf("Benchmark test item %d", i),
			"timestamp":   fmt.Sprintf("%d", i),
		}

		b.StartTimer()

		_, err := repo.Publish(streamName, data)

		b.StopTimer()

		if err != nil {
			b.Fatalf("Failed to publish to stream: %v", err)
		}
	}

	b.StopTimer()
	cleanupTestData(client, streamName)
}

func BenchmarkStreamRepositoryPublishParallel(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}

	client, err := setupTestRedis()
	if err != nil {
		b.Fatalf("Failed to setup test Redis: %v", err)
	}
	defer client.Close()

	repo := redis.NewStreamRepository(client)

	streamName := fmt.Sprintf("benchmark-stream-parallel-%s", uuid.New().String())

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			data := map[string]interface{}{
				"id":          uuid.New().String(),
				"description": fmt.Sprintf("Parallel benchmark test item %d", counter),
				"timestamp":   fmt.Sprintf("%d", counter),
			}
			counter++

			_, err := repo.Publish(streamName, data)
			if err != nil {
				b.Fatalf("Failed to publish to stream: %v", err)
			}
		}
	})

	b.StopTimer()
	cleanupTestData(client, streamName)
}
