package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server        ServerConfig        `json:"server"`
	Database      DatabaseConfig      `json:"database"`
	Redis         RedisConfig         `json:"redis"`
	CORS          CORSConfig          `json:"cors"`
	App           AppConfig           `json:"app"`
	ShopService   ShopServiceConfig   `json:"shop_service"`
	GRPC          GRPCConfig          `json:"grpc"`
	ShopServiceDB ShopServiceDBConfig `json:"shop_service_db"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host         string        `json:"host"`
	Port         string        `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
}

func (s *ServerConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host         string        `json:"host"`
	Port         string        `json:"port"`
	User         string        `json:"user"`
	Password     string        `json:"password"`
	DBName       string        `json:"db_name"`
	SSLMode      string        `json:"ssl_mode"`
	MaxOpenConns int           `json:"max_open_conns"`
	MaxIdleConns int           `json:"max_idle_conns"`
	MaxLifetime  time.Duration `json:"max_lifetime"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowOrigins []string `json:"allow_origins"`
	AllowMethods []string `json:"allow_methods"`
	AllowHeaders []string `json:"allow_headers"`
}

// AppConfig holds application configuration
type AppConfig struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
	LogLevel    string `json:"log_level"`
}

type ShopServiceConfig struct {
	Address string `json:"address"`
}

type GRPCConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type ShopServiceDBConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
	SSLMode  string `json:"ssl_mode"`
}

func (a *AppConfig) IsProduction() bool {
	return a.Environment == "production"
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file from project root
	_ = godotenv.Load("../../../../.env")

	config := &Config{
		Server: ServerConfig{
			Host:         getEnv("PRODUCT_SERVICE_SERVICE_HOST", "0.0.0.0"), // SỬA Ở ĐÂY
			Port:         getEnv("PRODUCT_SERVICE_SERVICE_PORT", "8082"),    // SỬA Ở ĐÂY
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Host:         getEnv("PRODUCT_SERVICE_POSTGRES_HOST", "localhost"),                // SỬA Ở ĐÂY
			Port:         getEnv("PRODUCT_SERVICE_POSTGRES_PORT_INTERNAL", "6002"),            // SỬA Ở ĐÂY
			User:         getEnv("PRODUCT_SERVICE_POSTGRES_USER", "postgres"),                 // SỬA Ở ĐÂY
			Password:     getEnv("PRODUCT_SERVICE_POSTGRES_PASSWORD", ""),                     // SỬA Ở ĐÂY
			DBName:       getEnv("PRODUCT_SERVICE_POSTGRES_DB", "product_service_go_shop_db"), // SỬA Ở ĐÂY
			SSLMode:      getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns: getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getIntEnv("DB_MAX_IDLE_CONNS", 25),
			MaxLifetime:  getDurationEnv("DB_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
		},
		CORS: CORSConfig{
			AllowOrigins: getStringSliceEnv("PRODUCT_SERVICE_CORS_ALLOWED_ORIGINS", []string{"*"}), // SỬA Ở ĐÂY
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"*"},
		},
		App: AppConfig{
			Name:        getEnv("PRODUCT_SERVICE_SERVICE_NAME", "product-service"), // SỬA Ở ĐÂY
			Version:     getEnv("APP_VERSION", "1.0.0"),
			Environment: getEnv("ENVIRONMENT", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
		},
		ShopService: ShopServiceConfig{
			Address: getEnv("SHOP_SERVICE_GRPC_ADDRESS", "shop-service:50051"), // SỬA Ở ĐÂY: Dùng biến môi trường mới và đúng service name
		},
		GRPC: GRPCConfig{
			Host: getEnv("PRODUCT_SERVICE_GRPC_HOST", "localhost"),
			Port: getEnv("PRODUCT_SERVICE_GRPC_PORT", "50051"),
		},
		ShopServiceDB: ShopServiceDBConfig{
			Host:     getEnv("SHOP_SERVICE_POSTGRES_HOST", "localhost"),
			Port:     getEnv("SHOP_SERVICE_POSTGRES_PORT", "6001"),
			User:     getEnv("SHOP_SERVICE_POSTGRES_USER", "postgres"),
			Password: getEnv("SHOP_SERVICE_POSTGRES_PASSWORD", "toai20102002"),
			DBName:   getEnv("SHOP_SERVICE_POSTGRES_DB", "shop_service_go_shop_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
	}

	return config, nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getStringSliceEnv(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Split by comma and trim spaces
		parts := strings.Split(value, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}

// GetDatabaseURL returns the database connection URL
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// GetServerAddress returns the server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

// GetRedisAddress returns the Redis address
func (c *Config) GetRedisAddress() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}
