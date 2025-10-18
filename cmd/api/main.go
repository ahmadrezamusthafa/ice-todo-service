package main

import (
	"fmt"
	"github.com/ahmadrezamusthafa/ice-todo-service/config"
	"github.com/ahmadrezamusthafa/ice-todo-service/infrastructure/database"
	"github.com/ahmadrezamusthafa/ice-todo-service/infrastructure/logger"
	"github.com/ahmadrezamusthafa/ice-todo-service/infrastructure/storage"
)

func main() {
	log := logger.NewStandardLogger()
	cfg := config.NewConfig()

	if err := cfg.ValidateConfig(); err != nil {
		log.Fatal("Configuration error: %v", err)
	}

	mysqlConnector := database.NewMySQLConnector(cfg, log)
	db, err := mysqlConnector.Connect()
	if err != nil {
		log.Fatal("Failed to connect to MySQL: %v", err)
	}
	defer mysqlConnector.Close()

	redisConnector := database.NewRedisConnector(cfg, log)
	redisClient, err := redisConnector.Connect()
	if err != nil {
		log.Fatal("Failed to connect to Redis: %v", err)
	}
	defer redisConnector.Close()

	s3Connector := storage.NewS3Connector(cfg, log)
	uploader, downloader, err := s3Connector.Connect()
	if err != nil {
		log.Fatal("Failed to connect to S3: %v", err)
	}

	fmt.Println(db, redisClient)
	fmt.Println(uploader, downloader)
}
