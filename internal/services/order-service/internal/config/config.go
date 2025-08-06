package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	App      AppConfig      `mapstructure:"app"`

	ShopServiceAdapter    ExternalServiceConfig `mapstructure:"shop_service_adapter"`
	ProductServiceAdapter ExternalServiceConfig `mapstructure:"product_service_adapter"`
	UserServiceAdapter    ExternalServiceConfig `mapstructure:"user_service_adapter"`
	GRPC                  GrpcConfig            `mapstructure:"grpc"`
	Kafka                 KafkaConfig           `mapstructure:"kafka"`
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
	ServiceHost string `mapstructure:"service_host"`
	ServicePort int    `mapstructure:"service_port"`
}

type ExternalServiceConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
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
			Host:         getEnv("ORDER_SERVICE_DB_HOST", "localhost"),
			Port:         getEnv("ORDER_SERVICE_POSTGRES_PORT_INTERNAL", "5432"),
			User:         getEnv("ORDER_SERVICE_DB_USER", "postgres"),
			Password:     getEnv("ORDER_SERVICE_DB_PASSWORD", "toai20102002"),
			DBName:       getEnv("ORDER_SERVICE_DB_NAME", "order_service_go_shop_db"),
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
		ShopServiceAdapter: ExternalServiceConfig{
			Host: getEnv("SHOP_SERVICE_GRPC_HOST", "localhost"),
			Port: getEnv("SHOP_SERVICE_GRPC_PORT", "8081"),
		},
		ProductServiceAdapter: ExternalServiceConfig{
			Host: getEnv("PRODUCT_SERVICE_GRPC_HOST", "localhost"),
			Port: getEnv("PRODUCT_SERVICE_GRPC_PORT", "8082"),
		},
		UserServiceAdapter: ExternalServiceConfig{
			Host: getEnv("USER_SERVICE_GRPC_HOST", "localhost"),
			Port: getEnv("USER_SERVICE_GRPC_PORT", "8084"),
		},
		GRPC: GrpcConfig{
			ServiceHost: getEnv("ORDER_SERVICE_GRPC_HOST", "localhost"),
			ServicePort: getIntEnv("ORDER_SERVICE_GRPC_PORT", 50052),
		},
		Kafka: KafkaConfig{
			Brokers: getSliceEnv("KAFKA_BROKERS", []string{"localhost:9092"}),
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

func getSliceEnv(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, ",")
}
