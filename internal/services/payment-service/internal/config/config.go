package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	App      AppConfig      `mapstructure:"app"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

func (s *ServerConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

type DatabaseConfig struct {
	Host         string        `mapstructure:"host"`
	Port         string        `mapstructure:"port"`
	User         string        `mapstructure:"user"`
	Password     string        `mapstructure:"password"`
	DBName       string        `mapstructure:"db_name"`
	SSLMode      string        `mapstructure:"ssl_mode"`
	MaxOpenConns int           `mapstructure:"max_open_conns"`
	MaxIdleConns int           `mapstructure:"max_idle_conns"`
	MaxLifetime  time.Duration `mapstructure:"max_lifetime"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
}

type GrpcConfig struct {
	ProductServiceHost string `mapstructure:"product_service_host"`
	ProductServicePort int    `mapstructure:"product_service_port"`
}

type ExternalServiceConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

func (a *AppConfig) IsProduction() bool {
	return a.Environment == "production"
}

func Load() (*Config, error) {
	cfg := &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "order-service"),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("ORDER_SERVICE_PORT", "8084"),
		},
		Database: DatabaseConfig{
			Host:         getEnv("PAYMENT_SERVICE_DB_HOST", "localhost"),
			Port:         getEnv("PAYMENT_SERVICE_POSTGRES_PORT_INTERNAL", "5432"),
			User:         getEnv("PAYMENT_SERVICE_DB_USER", "postgres"),
			Password:     getEnv("PAYMENT_SERVICE_DB_PASSWORD", "toai20102002"),
			DBName:       getEnv("PAYMENT_SERVICE_DB_NAME", "payment_service_go_shop_db"),
			SSLMode:      getEnv("DATABASE_SSLMODE", "disable"),
			MaxOpenConns: getIntEnv("DATABASE_MAX_OPEN_CONNS", 10),
			MaxIdleConns: getIntEnv("DATABASE_MAX_IDLE_CONNS", 5),
			MaxLifetime:  getDurationEnv("DATABASE_MAX_LIFETIME", time.Minute*5),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 1),
		},
	}
	return cfg, nil
}

// Các hàm helper để đọc biến môi trường
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}
