package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	S3       S3Config
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type S3Config struct {
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

func NewConfig() *Config {
	initViper()

	return &Config{
		Server: ServerConfig{
			Port: viper.GetString("PORT"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("MYSQL_HOST"),
			Port:     viper.GetString("MYSQL_PORT"),
			User:     viper.GetString("MYSQL_USER"),
			Password: viper.GetString("MYSQL_PASSWORD"),
			DBName:   viper.GetString("MYSQL_DATABASE"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
		},
		S3: S3Config{
			Endpoint:  viper.GetString("AWS_ENDPOINT"),
			Region:    viper.GetString("AWS_REGION"),
			AccessKey: viper.GetString("AWS_ACCESS_KEY"),
			SecretKey: viper.GetString("AWS_SECRET_KEY"),
			Bucket:    viper.GetString("S3_BUCKET"),
			UseSSL:    false,
		},
	}
}

func initViper() {

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("MYSQL_HOST", "localhost")
	viper.SetDefault("MYSQL_PORT", "3306")
	viper.SetDefault("MYSQL_USER", "root")
	viper.SetDefault("MYSQL_PASSWORD", "password")
	viper.SetDefault("MYSQL_DATABASE", "todo_db")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("AWS_ENDPOINT", "http://localhost:4566")
	viper.SetDefault("AWS_REGION", "us-east-1")
	viper.SetDefault("AWS_ACCESS_KEY", "test")
	viper.SetDefault("AWS_SECRET_KEY", "test")
	viper.SetDefault("S3_BUCKET", "todo-files")

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {

		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("Error reading config file: %s\n", err)
		}
	}

	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func (c *Config) ValidateConfig() error {

	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.Port == "" {
		return fmt.Errorf("database port is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	if c.S3.Region == "" {
		return fmt.Errorf("S3 region is required")
	}
	if c.S3.Bucket == "" {
		return fmt.Errorf("S3 bucket is required")
	}

	return nil
}
