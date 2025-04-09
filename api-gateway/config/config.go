package config

import (
	"github.com/spf13/viper"
)

// Config berisi semua konfigurasi aplikasi
type Config struct {
	Server    ServerConfig
	Services  ServicesConfig
	JWT       JWTConfig
	RateLimit RateLimitConfig
}

// ServerConfig untuk konfigurasi server
type ServerConfig struct {
	Port string
	Mode string // "debug", "release"
}

// ServicesConfig untuk URL service-service
type ServicesConfig struct {
	AuthURL         string
	UserURL         string
	FeedbackURL     string
	RewardURL       string
	InventoryURL    string
	ArticleURL      string
	NotificationURL string
}

// JWTConfig untuk konfigurasi JWT
type JWTConfig struct {
	Secret string
}

// RateLimitConfig untuk rate limiting
type RateLimitConfig struct {
	RequestsPerSecond int
}

// LoadConfig memuat konfigurasi dari environment atau file
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	// Default values
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("services.auth_url", "http://auth:8081")
	viper.SetDefault("services.user_url", "http://user:8082")
	viper.SetDefault("services.feedback_url", "http://localhost:8083")
	viper.SetDefault("services.reward_url", "http://reward:8084")
	viper.SetDefault("services.inventory_url", "http://inventory:8085")
	viper.SetDefault("services.article_url", "http://article:8086")
	viper.SetDefault("services.notification_url", "http://notification:8087")
	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("rate_limit.requests_per_second", 100)

	// Environment variables mapping
	viper.BindEnv("server.port", "PORT")
	viper.BindEnv("server.mode", "GIN_MODE")
	viper.BindEnv("services.auth_url", "AUTH_SERVICE_URL")
	viper.BindEnv("services.user_url", "USER_SERVICE_URL")
	viper.BindEnv("services.feedback_url", "FEEDBACK_SERVICE_URL")
	viper.BindEnv("services.reward_url", "REWARD_SERVICE_URL")
	viper.BindEnv("services.inventory_url", "INVENTORY_SERVICE_URL")
	viper.BindEnv("services.article_url", "ARTICLE_SERVICE_URL")
	viper.BindEnv("services.notification_url", "NOTIFICATION_SERVICE_URL")
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("rate_limit.requests_per_second", "RATE_LIMIT_RPS")

	// Try to read config file if it exists
	err := viper.ReadInConfig()
	if err != nil {
		// It's ok if config file doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	config.Server.Port = viper.GetString("server.port")
	config.Server.Mode = viper.GetString("server.mode")
	config.Services.AuthURL = viper.GetString("services.auth_url")
	config.Services.UserURL = viper.GetString("services.user_url")
	config.Services.FeedbackURL = viper.GetString("services.feedback_url")
	config.Services.RewardURL = viper.GetString("services.reward_url")
	config.Services.InventoryURL = viper.GetString("services.inventory_url")
	config.Services.ArticleURL = viper.GetString("services.article_url")
	config.Services.NotificationURL = viper.GetString("services.notification_url")
	config.JWT.Secret = viper.GetString("jwt.secret")
	config.RateLimit.RequestsPerSecond = viper.GetInt("rate_limit.requests_per_second")

	return &config, nil
}
