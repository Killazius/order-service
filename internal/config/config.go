package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"order-service/internal/storage/postgresql"
)

type Config struct {
	GPRCport int               `yaml:"GRPC_PORT" env:"GRPC_PORT" env-default:"50051"`
	HTTPport int               `yaml:"HTTP_PORT" env:"HTTP_PORT" env-default:"8080"`
	Postgres postgresql.Config `yaml:"POSTGRES"`
}

func MustLoad() (*Config, error) {
	_ = godotenv.Load()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("CONFIG_PATH does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	if postgresHost := os.Getenv("POSTGRES_HOST"); postgresHost != "" {
		cfg.Postgres.Host = postgresHost
	}
	if postgresPort := os.Getenv("POSTGRES_PORT"); postgresPort != "" {
		port, err := strconv.Atoi(postgresPort)
		if err != nil {
			return nil, fmt.Errorf("invalid POSTGRES_PORT: %w", err)
		}
		cfg.Postgres.Port = uint16(port)
	}
	if postgresUser := os.Getenv("POSTGRES_USER"); postgresUser != "" {
		cfg.Postgres.Username = postgresUser
	}
	if postgresPassword := os.Getenv("POSTGRES_PASSWORD"); postgresPassword != "" {
		cfg.Postgres.Password = postgresPassword
	}
	if postgresDB := os.Getenv("POSTGRES_DB"); postgresDB != "" {
		cfg.Postgres.Database = postgresDB
	}
	if postgresMaxConn := os.Getenv("POSTGRES_MAX_CONN"); postgresMaxConn != "" {
		maxConn, err := strconv.Atoi(postgresMaxConn)
		if err != nil {
			return nil, fmt.Errorf("invalid POSTGRES_MAX_CONN: %w", err)
		}
		cfg.Postgres.MaxConn = int32(maxConn)
	}
	if postgresMinConn := os.Getenv("POSTGRES_MIN_CONN"); postgresMinConn != "" {
		minConn, err := strconv.Atoi(postgresMinConn)
		if err != nil {
			return nil, fmt.Errorf("invalid POSTGRES_MIN_CONN: %w", err)
		}
		cfg.Postgres.MinConn = int32(minConn)
	}
	return &cfg, nil
}
