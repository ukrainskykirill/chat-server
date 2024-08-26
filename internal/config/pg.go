package config

import (
	"fmt"
	"os"
)

const (
	pgHostEnv = "DB_HOST"
	pgPortEnv = "DB_PORT"
	pgUserEnv = "DB_USER"
	pgPassEnv = "DB_PASSWORD"
	pgDBEnv   = "DB_DATABASE_NAME"
)

type PGConfig interface {
	URL() string
}

type pgConfig struct {
	host     string
	port     string
	username string
	password string
	database string
}

func NewDBConfig() (PGConfig, error) {
	host, ok := os.LookupEnv(pgHostEnv)
	if !ok {
		return &pgConfig{}, errVariableNotFound
	}

	port, ok := os.LookupEnv(pgPortEnv)
	if !ok {
		return &pgConfig{}, errVariableNotFound
	}

	username, ok := os.LookupEnv(pgUserEnv)
	if !ok {
		return &pgConfig{}, errVariableNotFound
	}

	password, ok := os.LookupEnv(pgPassEnv)
	if !ok {
		return &pgConfig{}, errVariableNotFound
	}

	database, ok := os.LookupEnv(pgDBEnv)
	if !ok {
		return &pgConfig{}, errVariableNotFound
	}

	return &pgConfig{
		host:     host,
		port:     port,
		username: username,
		password: password,
		database: database,
	}, nil
}

func (c *pgConfig) URL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.username, c.password, c.host, c.port, c.database)
}
