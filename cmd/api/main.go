package main

import (
	"github.com/ahmadrezamusthafa/ice-todo-service/config"
	"github.com/ahmadrezamusthafa/ice-todo-service/infrastructure/logger"
)

func main() {
	log := logger.NewStandardLogger()
	cfg := config.NewConfig()

	if err := cfg.ValidateConfig(); err != nil {
		log.Fatal("Configuration error: %v", err)
	}
}
