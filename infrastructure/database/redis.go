package database

import (
	"context"
	"fmt"
	"github.com/ahmadrezamusthafa/ice-todo-service/config"
	"github.com/ahmadrezamusthafa/ice-todo-service/infrastructure/logger"

	"github.com/go-redis/redis/v8"
)

type RedisConnector struct {
	config *config.Config
	logger logger.Logger
	client *redis.Client
}

func NewRedisConnector(config *config.Config, logger logger.Logger) *RedisConnector {
	return &RedisConnector{
		config: config,
		logger: logger,
	}
}

func (c *RedisConnector) Connect() (*redis.Client, error) {
	if c.client != nil {
		return c.client, nil
	}

	c.client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", c.config.Redis.Host, c.config.Redis.Port),
		Password: c.config.Redis.Password,
		DB:       0,
	})

	ctx := context.Background()
	pong, err := c.client.Ping(ctx).Result()
	if err != nil {
		c.logger.Error("Failed to connect to Redis: %v", err)
		return nil, err
	}

	c.logger.Info("Connected to Redis: %s", pong)

	return c.client, nil
}

func (c *RedisConnector) Close() error {
	if c.client != nil {
		c.logger.Info("Closing Redis connection")
		return c.client.Close()
	}

	return nil
}
