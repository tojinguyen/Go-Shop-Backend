package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the user service
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	JWT      JWTConfig      `json:"jwt"`
	Redis    RedisConfig    `json:"redis"`
	CORS     CORSConfig     `json:"cors"`
	App      AppConfig      `json:"app"`
	Email    EmailConfig    `json:"email"`
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
	Name         string        `json:"name"`
	SSLMode      string        `json:"ssl_mode"`
	MaxOpenConns int           `json:"max_open_conns"`
	MaxIdleConns int           `json:"max_idle_conns"`
	MaxLifetime  time.Duration `json:"max_lifetime"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey       string        `json:"secret_key"`
	AccessTokenTTL  time.Duration `json:"access_token_ttl"`
	RefreshTokenTTL time.Duration `json:"refresh_token_ttl"`
	Issuer          string        `json:"issuer"`
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
	AllowedOrigins   []string `json:"allowed_origins"`
	AllowedMethods   []string `json:"allowed_methods"`
	AllowedHeaders   []string `json:"allowed_headers"`
	AllowCredentials bool     `json:"allow_credentials"`
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
	Debug       bool   `json:"debug"`
	LogLevel    string `json:"log_level"`
	FrontendURL string `json:"frontend_url"`
}

// EmailConfig holds email configuration
type EmailConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	From         string `json:"from"`
	FromName     string `json:"from_name"`
	UseTLS       bool   `json:"use_tls"`
	UseSSL       bool   `json:"use_ssl"`
	TemplatePath string `json:"template_path"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Host:         getEnvWithDefault("SERVER_HOST", "localhost"),
			Port:         getEnvWithDefault("SERVER_PORT", "8081"),
			ReadTimeout:  getDurationEnvWithDefault("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnvWithDefault("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getDurationEnvWithDefault("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Host:         getEnvWithDefault("DB_HOST", "localhost"),
			Port:         getEnvWithDefault("DB_PORT", "5432"),
			User:         getEnvWithDefault("DB_USER", "postgres"),
			Password:     getEnvWithDefault("DB_PASSWORD", ""),
			Name:         getEnvWithDefault("DB_NAME", "go_shop_users"),
			SSLMode:      getEnvWithDefault("DB_SSL_MODE", "disable"),
			MaxOpenConns: getIntEnvWithDefault("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getIntEnvWithDefault("DB_MAX_IDLE_CONNS", 25),
			MaxLifetime:  getDurationEnvWithDefault("DB_MAX_LIFETIME", 5*time.Minute),
		},
		JWT: JWTConfig{
			SecretKey:       getEnvWithDefault("JWT_SECRET_KEY", "your-secret-key"),
			AccessTokenTTL:  getDurationEnvWithDefault("JWT_ACCESS_TOKEN_TTL", 15*time.Minute),
			RefreshTokenTTL: getDurationEnvWithDefault("JWT_REFRESH_TOKEN_TTL", 24*time.Hour),
			Issuer:          getEnvWithDefault("JWT_ISSUER", "go-shop-user-service"),
		},
		Redis: RedisConfig{
			Host:     getEnvWithDefault("REDIS_HOST", "localhost"),
			Port:     getEnvWithDefault("REDIS_PORT", "6379"),
			Password: getEnvWithDefault("REDIS_PASSWORD", ""),
			DB:       getIntEnvWithDefault("REDIS_DB", 0),
		},
		CORS: CORSConfig{
			AllowedOrigins:   getStringSliceEnvWithDefault("CORS_ALLOWED_ORIGINS", []string{"*"}),
			AllowedMethods:   getStringSliceEnvWithDefault("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowedHeaders:   getStringSliceEnvWithDefault("CORS_ALLOWED_HEADERS", []string{"*"}),
			AllowCredentials: getBoolEnvWithDefault("CORS_ALLOW_CREDENTIALS", true),
		},
		App: AppConfig{
			Name:        getEnvWithDefault("APP_NAME", "Go-Shop User Service"),
			Version:     getEnvWithDefault("APP_VERSION", "1.0.0"),
			Environment: getEnvWithDefault("APP_ENV", "development"),
			Debug:       getBoolEnvWithDefault("APP_DEBUG", true),
			LogLevel:    getEnvWithDefault("LOG_LEVEL", "info"),
			FrontendURL: getEnvWithDefault("FRONTEND_URL", "http://localhost:3000"),
		},
		Email: EmailConfig{
			Host:         getEnvWithDefault("EMAIL_HOST", "smtp.example.com"),
			Port:         getIntEnvWithDefault("EMAIL_PORT", 587),
			Username:     getEnvWithDefault("EMAIL_USERNAME", ""),
			Password:     getEnvWithDefault("EMAIL_PASSWORD", ""),
			From:         getEnvWithDefault("EMAIL_FROM", ""),
			FromName:     getEnvWithDefault("EMAIL_FROM_NAME", "Go-Shop"),
			UseTLS:       getBoolEnvWithDefault("EMAIL_USE_TLS", true),
			UseSSL:       getBoolEnvWithDefault("EMAIL_USE_SSL", false),
			TemplatePath: getEnvWithDefault("EMAIL_TEMPLATE_PATH", "./templates/email"),
		},
	}

	return config, nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// GetServerAddress returns the server address
func (c *ServerConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// GetRedisAddress returns the Redis address
func (c *RedisConfig) GetRedisAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// IsProduction returns true if the environment is production
func (c *AppConfig) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment returns true if the environment is development
func (c *AppConfig) IsDevelopment() bool {
	return c.Environment == "development"
}

// Helper functions for environment variable parsing
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnvWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnvWithDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnvWithDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getStringSliceEnvWithDefault(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Simple split by comma, you might want to use a more sophisticated parser
		return []string{value}
	}
	return defaultValue
}
