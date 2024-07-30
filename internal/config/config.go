package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	errVariableNotFound = errors.New("environment variable not found")
)

type dbConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	URL      string
}

func loadDBConfig() (*dbConfig, error) {
	host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		return &dbConfig{}, errVariableNotFound
	}

	port, ok := os.LookupEnv("DB_PORT")
	if !ok {
		return &dbConfig{}, errVariableNotFound
	}
	intPort, err := strconv.Atoi(port)
	if err != nil {
		return &dbConfig{}, err
	}

	username, ok := os.LookupEnv("DB_USER")
	if !ok {
		return &dbConfig{}, errVariableNotFound
	}

	password, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		return &dbConfig{}, errVariableNotFound
	}

	database, ok := os.LookupEnv("DB_DATABASE_NAME")
	if !ok {
		return &dbConfig{}, errVariableNotFound
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, host, intPort, database)
	return &dbConfig{
		Host:     host,
		Port:     intPort,
		Username: username,
		Password: password,
		Database: database,
		URL:      url,
	}, nil
}

type grpcConfig struct {
	Port int
}

func loadGRPCConfig() (*grpcConfig, error) {
	port, ok := os.LookupEnv("GRPC_PORT")
	if !ok {
		return &grpcConfig{}, errVariableNotFound
	}

	intPort, err := strconv.Atoi(port)
	if err != nil {
		return &grpcConfig{}, err
	}
	return &grpcConfig{
		Port: intPort,
	}, nil
}

// AppConfig app config struct
type AppConfig struct {
	GRPC *grpcConfig
	DB   *dbConfig
}

func loadAppConfig() (*AppConfig, error) {
	dbConfig, err := loadDBConfig()
	if err != nil {
		return &AppConfig{}, err
	}

	grpcConfig, err := loadGRPCConfig()
	if err != nil {
		return &AppConfig{}, err
	}
	return &AppConfig{
		GRPC: grpcConfig,
		DB:   dbConfig,
	}, err
}

// InitConfig init app config
func InitConfig() (*AppConfig, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	config, err := loadAppConfig()
	if err != nil {
		return &AppConfig{}, err
	}
	return config, nil
}
