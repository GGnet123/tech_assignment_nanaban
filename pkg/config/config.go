package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppName    string
	AppEnv     string
	Server     ServerConfig
	DB         DBConfig
	Prometheus PrometheusConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type PrometheusConfig struct {
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func Load() (*Config, error) {
	cfg := &Config{
		AppName: getEnv("APP_NAME", "nanaban"),
		AppEnv:  getEnv("APP_ENV", "dev"),
		Prometheus: PrometheusConfig{
			Port: getEnv("PROMETHEUS_PORT", "9090"),
		},
		Server: ServerConfig{
			Host: getEnv("APP_HOST", "http://localhost"),
			Port: getEnv("APP_PORT", "8080"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "postgres"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "exchange"),
			Password: getEnv("DB_PASSWORD", "exchange_password"),
			DBName:   getEnv("DB_NAME", "exchange_db"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		c.DB.Host,
		c.DB.Port,
		c.DB.User,
		c.DB.Password,
		c.DB.DBName,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("APP_PORT is required")
	}
	if c.DB.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.DB.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.DB.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	return nil
}
