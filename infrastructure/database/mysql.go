package database

import (
	"database/sql"
	"fmt"
	"github.com/ahmadrezamusthafa/ice-todo-service/config"
	"github.com/ahmadrezamusthafa/ice-todo-service/infrastructure/logger"

	"github.com/go-sql-driver/mysql"
)

type MySQLConnector struct {
	config *config.Config
	logger logger.Logger
	db     *sql.DB
}

func NewMySQLConnector(config *config.Config, logger logger.Logger) *MySQLConnector {
	return &MySQLConnector{
		config: config,
		logger: logger,
	}
}

func (c *MySQLConnector) Connect() (*sql.DB, error) {
	if c.db != nil {
		return c.db, nil
	}

	mysqlConfig := mysql.Config{
		User:                 c.config.Database.User,
		Passwd:               c.config.Database.Password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", c.config.Database.Host, c.config.Database.Port),
		DBName:               c.config.Database.DBName,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err := sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		c.logger.Error("Failed to connect to MySQL: %v", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		c.logger.Error("Failed to ping MySQL: %v", err)
		return nil, err
	}

	c.db = db
	c.logger.Info("Connected to MySQL database")

	return c.db, nil
}

func (c *MySQLConnector) Close() error {
	if c.db != nil {
		c.logger.Info("Closing MySQL connection")
		return c.db.Close()
	}

	return nil
}
