package main

import (
	"fmt"
	"github.com/ahmadrezamusthafa/ice-todo-service/config"
	"github.com/ahmadrezamusthafa/ice-todo-service/infrastructure/database"
	"github.com/ahmadrezamusthafa/ice-todo-service/infrastructure/logger"
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

	fmt.Println(db, redisClient)
}
