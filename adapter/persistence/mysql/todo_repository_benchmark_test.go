package mysql_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/ahmadrezamusthafa/ice-todo-service/adapter/persistence/mysql"
	"github.com/ahmadrezamusthafa/ice-todo-service/config"
	"github.com/ahmadrezamusthafa/ice-todo-service/domain/entity"
	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func setupTestDB() (*sql.DB, error) {
	cfg := config.NewConfig()

	mysqlConfig := mysqldriver.Config{
		User:                 cfg.Database.User,
		Passwd:               cfg.Database.Password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", cfg.Database.Host, cfg.Database.Port),
		DBName:               cfg.Database.DBName,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err := sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping MySQL: %w", err)
	}

	return db, nil
}

func cleanupTestData(db *sql.DB, ids []string) {
	for _, id := range ids {
		_, _ = db.Exec("DELETE FROM todo_items WHERE id = ?", id)
	}
}

func BenchmarkTodoRepositoryCreate(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}

	db, err := setupTestDB()
	if err != nil {
		b.Fatalf("Failed to setup test database: %v", err)
	}
	defer db.Close()

	repo := mysql.NewTodoRepository(db)

	var createdIDs []string

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		id := uuid.New()
		todo := &entity.TodoItem{
			ID:          id,
			Description: fmt.Sprintf("Benchmark test item %d", i),
			DueDate:     time.Now().Add(24 * time.Hour),
			FileID:      "",
		}

		createdIDs = append(createdIDs, id.String())

		b.StartTimer()

		_, err := repo.Create(todo)

		b.StopTimer()

		if err != nil {
			b.Fatalf("Failed to create todo item: %v", err)
		}
	}

	b.StopTimer()
	cleanupTestData(db, createdIDs)
}

func BenchmarkTodoRepositoryCreateParallel(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}

	db, err := setupTestDB()
	if err != nil {
		b.Fatalf("Failed to setup test database: %v", err)
	}
	defer db.Close()

	repo := mysql.NewTodoRepository(db)

	idChan := make(chan string, b.N)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {

			id := uuid.New()
			todo := &entity.TodoItem{
				ID:          id,
				Description: fmt.Sprintf("Parallel benchmark test item %s", id),
				DueDate:     time.Now().Add(24 * time.Hour),
				FileID:      "",
			}

			_, err := repo.Create(todo)
			if err != nil {
				b.Fatalf("Failed to create todo item: %v", err)
			}

			idChan <- id.String()
		}
	})

	b.StopTimer()
	close(idChan)
	var createdIDs []string
	for id := range idChan {
		createdIDs = append(createdIDs, id)
	}

	cleanupTestData(db, createdIDs)
}
