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
	Server    ServerConfig    `json:"server"`
	Database  DatabaseConfig  `json:"database"`
	JWT       JWTConfig       `json:"jwt"`
	Redis     RedisConfig     `json:"redis"`
	CORS      CORSConfig      `json:"cors"`
	App       AppConfig       `json:"app"`
	Email     EmailConfig     `json:"email"`
	GRPC      GRPCConfig      `json:"grpc"`
	RateLimit RateLimitConfig `json:"rate_limit"`
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

type GRPCConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type RateLimitConfig struct {
	LoginMaxAttempts   int           `json:"login_max_attempts"`
	LoginWindowMinutes time.Duration `json:"login_window_minutes"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Host:         getEnvWithDefault("USER_SERVICE_SERVICE_HOST", "0.0.0.0"),
			Port:         getEnvWithDefault("USER_SERVICE_SERVICE_PORT", "8080"),
			ReadTimeout:  getDurationEnvWithDefault("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnvWithDefault("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getDurationEnvWithDefault("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Host:         getEnvWithDefault("USER_SERVICE_POSTGRES_HOST", "localhost"),
			Port:         getEnvWithDefault("USER_SERVICE_POSTGRES_PORT_INTERNAL", "5432"),
			User:         getEnvWithDefault("USER_SERVICE_POSTGRES_USER", "postgres"),
			Password:     getEnvWithDefault("USER_SERVICE_POSTGRES_PASSWORD", "postgres123"),
			Name:         getEnvWithDefault("USER_SERVICE_POSTGRES_DB", "go_shop_user_service"),
			SSLMode:      getEnvWithDefault("DB_SSL_MODE", "disable"),
			MaxOpenConns: getIntEnvWithDefault("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getIntEnvWithDefault("DB_MAX_IDLE_CONNS", 25),
			MaxLifetime:  getDurationEnvWithDefault("DB_MAX_LIFETIME", 5*time.Minute),
		},
		JWT: JWTConfig{
			SecretKey:       getEnvWithDefault("JWT_SECRET_KEY", "your-secret-key"),
			AccessTokenTTL:  getDurationEnvWithDefault("JWT_ACCESS_TOKEN_EXPIRY", 15*time.Minute),
			RefreshTokenTTL: getDurationEnvWithDefault("JWT_REFRESH_TOKEN_EXPIRY", 24*time.Hour),
			Issuer:          getEnvWithDefault("JWT_ISSUER", "go-shop-user-service"),
		},
		Redis: RedisConfig{
			Host:     getEnvWithDefault("REDIS_HOST", "redis-cache"),
			Port:     getEnvWithDefault("REDIS_PORT", "6379"),
			Password: getEnvWithDefault("REDIS_PASSWORD", "redis123"),
			DB:       getIntEnvWithDefault("REDIS_DB", 0),
		},
		CORS: CORSConfig{
			AllowedOrigins:   getStringSliceEnvWithDefault("USER_SERVICE_CORS_ALLOWED_ORIGINS", []string{"*"}),
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
			Host:         getEnvWithDefault("USER_SERVICE_SMTP_HOST", "smtp.example.com"),
			Port:         getIntEnvWithDefault("USER_SERVICE_SMTP_PORT", 587),
			Username:     getEnvWithDefault("USER_SERVICE_SMTP_USER", ""),
			Password:     getEnvWithDefault("USER_SERVICE_SMTP_PASSWORD", ""),
			From:         getEnvWithDefault("USER_SERVICE_EMAIL_FROM", ""),
			FromName:     getEnvWithDefault("USER_SERVICE_EMAIL_FROM_NAME", "Go-Shop"),
			UseTLS:       getBoolEnvWithDefault("USER_SERVICE_SMTP_USE_TLS", true),
			UseSSL:       getBoolEnvWithDefault("USER_SERVICE_SMTP_USE_SSL", false),
			TemplatePath: getEnvWithDefault("USER_SERVICE_EMAIL_TEMPLATE_PATH", "./templates/email"),
		},
		GRPC: GRPCConfig{
			Host: getEnvWithDefault("USER_SERVICE_GRPC_HOST", "localhost"),
			Port: getEnvWithDefault("USER_SERVICE_GRPC_PORT", "50051"),
		},
		RateLimit: RateLimitConfig{
			LoginMaxAttempts:   getIntEnvWithDefault("RATE_LIMIT_LOGIN_ATTEMPTS", 5),
			LoginWindowMinutes: getDurationEnvWithDefault("RATE_LIMIT_LOGIN_WINDOW_MINUTES", 1*time.Minute),
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

// Debug prints configuration values (without sensitive data)
func (c *Config) Debug() {
	fmt.Printf("=== Configuration Debug ===\n")
	fmt.Printf("Server: %s:%s\n", c.Server.Host, c.Server.Port)
	fmt.Printf("Database: %s@%s:%s/%s (SSL: %s)\n", c.Database.User, c.Database.Host, c.Database.Port, c.Database.Name, c.Database.SSLMode)
	fmt.Printf("Database Password: %s\n", c.Database.Password)
	fmt.Printf("Redis: %s:%s (DB: %d)\n", c.Redis.Host, c.Redis.Port, c.Redis.DB)
	fmt.Printf("JWT Issuer: %s\n", c.JWT.Issuer)
	fmt.Printf("Environment: %s\n", c.App.Environment)
	fmt.Printf("Debug Mode: %v\n", c.App.Debug)
	fmt.Printf("===========================\n")
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
