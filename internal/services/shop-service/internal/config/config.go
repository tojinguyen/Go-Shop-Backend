package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the shop service
type Config struct {
	Server        ServerConfig        `json:"server"`
	Database      DatabaseConfig      `json:"database"`
	JWT           JWTConfig           `json:"jwt"`
	Redis         RedisConfig         `json:"redis"`
	CORS          CORSConfig          `json:"cors"`
	App           AppConfig           `json:"app"`
	GRPC          GRPCConfig          `json:"grpc"`
	UserServiceDB UserServiceDBConfig `json:"user_service_db"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host         string        `json:"host"`
	Port         string        `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
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

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey            string        `json:"secret_key"`
	AccessTokenDuration  time.Duration `json:"access_token_duration"`
	RefreshTokenDuration time.Duration `json:"refresh_token_duration"`
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

type GRPCConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type UserServiceDBConfig struct {
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
			Host:         getEnv("SHOP_SERVICE_SERVICE_HOST", "0.0.0.0"),
			Port:         getEnv("SHOP_SERVICE_SERVICE_PORT", "8081"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Host:         getEnv("SHOP_SERVICE_POSTGRES_HOST", "localhost"),
			Port:         getEnv("SHOP_SERVICE_POSTGRES_PORT_INTERNAL", "6001"),
			User:         getEnv("SHOP_SERVICE_POSTGRES_USER", "postgres"),
			Password:     getEnv("SHOP_SERVICE_POSTGRES_PASSWORD", ""),
			DBName:       getEnv("SHOP_SERVICE_POSTGRES_DB", "shop_service_go_shop_db"),
			SSLMode:      getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns: getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getIntEnv("DB_MAX_IDLE_CONNS", 25),
			MaxLifetime:  getDurationEnv("DB_MAX_LIFETIME", 5*time.Minute),
		},
		JWT: JWTConfig{
			SecretKey:            getEnv("JWT_SECRET_KEY", "your-secret-key"),
			AccessTokenDuration:  getDurationEnv("JWT_ACCESS_TOKEN_EXPIRY", 24*time.Hour),
			RefreshTokenDuration: getDurationEnv("JWT_REFRESH_TOKEN_EXPIRY", 168*time.Hour),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
		},
		CORS: CORSConfig{
			AllowOrigins: getStringSliceEnv("SHOP_SERVICE_CORS_ALLOWED_ORIGINS", []string{"*"}),
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"*"},
		},
		App: AppConfig{
			Name:        getEnv("SHOP_SERVICE_SERVICE_NAME", "shop-service"),
			Version:     getEnv("APP_VERSION", "1.0.0"),
			Environment: getEnv("ENVIRONMENT", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
		},
		GRPC: GRPCConfig{
			Host: getEnv("SHOP_SERVICE_GRPC_HOST", "0.0.0.0"),
			Port: getEnv("SHOP_SERVICE_GRPC_PORT", "50051"),
		},
		UserServiceDB: UserServiceDBConfig{
			Host:     getEnv("USER_SERVICE_POSTGRES_HOST", "localhost"),
			Port:     getEnv("USER_SERVICE_POSTGRES_PORT_INTERNAL", "6000"), // Port public
			User:     getEnv("USER_SERVICE_POSTGRES_USER", "postgres"),
			Password: getEnv("USER_SERVICE_POSTGRES_PASSWORD", "toai20102002"),
			DBName:   getEnv("USER_SERVICE_POSTGRES_DB", "user_service_go_shop_db"),
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
