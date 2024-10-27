package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Logger    LoggerConfig
	OrmConfig OrmConfig
}

type ServerConfig struct {
	Port    string
	Timeout int
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type LoggerConfig struct {
	Level string
	File  string
}

type OrmConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	QueryLog        bool
}

func Load() (*Config, error) {
	err := godotenv.Load()

	if err != nil {
		return nil, err
	}

	timeout, err := strconv.Atoi(os.Getenv("SERVER_TIMEOUT"))
	if err != nil {
		timeout = 30 // default value
	}

	serverConfig := ServerConfig{
		Port:    os.Getenv("PORT"),
		Timeout: timeout,
	}

	databaseConfig := DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	loggerConfig := LoggerConfig{
		Level: os.Getenv("LOGGER_LEVEL"),
		File:  os.Getenv("LOGGER_FILE"),
	}

	ormConfig := OrmConfig{
		MaxOpenConns:    20,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
		QueryLog:        true,
	}

	config := Config{
		Server:    serverConfig,
		Database:  databaseConfig,
		Logger:    loggerConfig,
		OrmConfig: ormConfig,
	}

	return &config, nil
}
