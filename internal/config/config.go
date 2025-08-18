package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string `mapstructure:"environment"`
	HTTPPort    string `mapstructure:"http_port"`
	MetricsPort string `mapstructure:"metrics_port"`

	Database   DatabaseConfig   `mapstructure:"database"`
	ClickHouse ClickHouseConfig `mapstructure:"clickhouse"`
	Redis      RedisConfig      `mapstructure:"redis"`
	NATS       NATSConfig       `mapstructure:"nats"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Security   SecurityConfig   `mapstructure:"security"`
	RateLimit  RateLimitConfig  `mapstructure:"rate_limit"`
	AI         AIConfig         `mapstructure:"ai"`

	// Service-specific configurations (to be customized per MCP)
	{{SERVICE_CONFIG_NAME}} {{SERVICE_CONFIG_TYPE}} `mapstructure:"{{SERVICE_CONFIG_KEY}}"`
}

type DatabaseConfig struct {
	URL             string `mapstructure:"url"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime string `mapstructure:"conn_max_lifetime"`
}

type ClickHouseConfig struct {
	URL      string `mapstructure:"url"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type RedisConfig struct {
	URL        string `mapstructure:"url"`
	DB         int    `mapstructure:"db"`
	MaxRetries int    `mapstructure:"max_retries"`
}

type NATSConfig struct {
	URL     string `mapstructure:"url"`
	Cluster string `mapstructure:"cluster"`
}

type JWTConfig struct {
	Secret   string `mapstructure:"secret"`
	Issuer   string `mapstructure:"issuer"`
	Audience string `mapstructure:"audience"`
}

type SecurityConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	APIKey         string   `mapstructure:"api_key"`
}

type RateLimitConfig struct {
	Enabled bool `mapstructure:"enabled"`
	RPS     int  `mapstructure:"rps"`
	Burst   int  `mapstructure:"burst"`
}

type AIConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Provider string `mapstructure:"provider"`
	APIKey   string `mapstructure:"api_key"`
	Model    string `mapstructure:"model"`
	BaseURL  string `mapstructure:"base_url"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Set defaults
	setDefaults()

	// Read environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Override with environment variables
	overrideWithEnv(&config)

	return &config, nil
}

func setDefaults() {
	viper.SetDefault("environment", "development")
	viper.SetDefault("http_port", "8080")
	viper.SetDefault("metrics_port", "9090")

	// Database defaults
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 25)
	viper.SetDefault("database.conn_max_lifetime", "5m")

	// Redis defaults
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.max_retries", 3)

	// Security defaults
	viper.SetDefault("security.allowed_origins", []string{"http://localhost:3000"})

	// Rate limit defaults
	viper.SetDefault("rate_limit.enabled", true)
	viper.SetDefault("rate_limit.rps", 100)
	viper.SetDefault("rate_limit.burst", 200)

	// AI defaults
	viper.SetDefault("ai.enabled", true)
	viper.SetDefault("ai.provider", "openai")
	viper.SetDefault("ai.model", "gpt-4")
}

func overrideWithEnv(config *Config) {
	if port := os.Getenv("PORT"); port != "" {
		config.HTTPPort = port
	}

	if env := os.Getenv("ENVIRONMENT"); env != "" {
		config.Environment = env
	}

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		config.Database.URL = dbURL
	}

	if clickhouseURL := os.Getenv("CLICKHOUSE_URL"); clickhouseURL != "" {
		config.ClickHouse.URL = clickhouseURL
	}

	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		config.Redis.URL = redisURL
	}

	if natsURL := os.Getenv("NATS_URL"); natsURL != "" {
		config.NATS.URL = natsURL
	}

	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		config.JWT.Secret = jwtSecret
	}

	if apiKey := os.Getenv("API_KEY"); apiKey != "" {
		config.Security.APIKey = apiKey
	}

	if aiAPIKey := os.Getenv("AI_API_KEY"); aiAPIKey != "" {
		config.AI.APIKey = aiAPIKey
	}

	if rpsStr := os.Getenv("RATE_LIMIT_RPS"); rpsStr != "" {
		if rps, err := strconv.Atoi(rpsStr); err == nil {
			config.RateLimit.RPS = rps
		}
	}
}
