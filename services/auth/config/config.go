package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config berisi semua konfigurasi aplikasi
type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	JWT        JWTConfig
	RabbitMQ   RabbitMQConfig
	SMTP       SMTPConfig
	GoogleAuth GoogleAuthConfig
}

// ServerConfig untuk konfigurasi server
type ServerConfig struct {
	Port string
	Mode string // "debug", "release"
}

// DatabaseConfig untuk konfigurasi database
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// JWTConfig untuk konfigurasi JWT
type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

// RabbitMQConfig untuk konfigurasi RabbitMQ
type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Exchange string
}

// SMTPConfig untuk konfigurasi SMTP
type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
}

// GoogleAuthConfig untuk konfigurasi Google OAuth
type GoogleAuthConfig struct {
	WebClientID     string
	AndroidClientID string
	IOSClientID     string
	RedirectURL     string
}

// LoadConfig memuat konfigurasi dari environment atau file
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	// Default values
	viper.SetDefault("server.port", "8081")
	viper.SetDefault("server.mode", "debug")

	viper.SetDefault("database.host", "auth-db")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.name", "auth_db")
	viper.SetDefault("database.sslmode", "disable")

	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.expiry", "24h")

	viper.SetDefault("rabbitmq.host", "rabbitmq")
	viper.SetDefault("rabbitmq.port", "5672")
	viper.SetDefault("rabbitmq.user", "guest")
	viper.SetDefault("rabbitmq.password", "guest")
	viper.SetDefault("rabbitmq.exchange", "auth.events")

	viper.SetDefault("smtp.host", "smtp.gmail.com")
	viper.SetDefault("smtp.port", "587")
	viper.SetDefault("smtp.user", "")
	viper.SetDefault("smtp.password", "")
	viper.SetDefault("smtp.from", "no-reply@lrt-crm.com")

	viper.SetDefault("oauth.google.web_client_id", "")
	viper.SetDefault("oauth.google.android_client_id", "")
	viper.SetDefault("oauth.google.ios_client_id", "")
	viper.SetDefault("oauth.google.redirect_url", "")

	// Environment variables mapping
	viper.BindEnv("server.port", "PORT")
	viper.BindEnv("server.mode", "GIN_MODE")

	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.name", "DB_NAME")
	viper.BindEnv("database.sslmode", "DB_SSLMODE")

	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("jwt.expiry", "JWT_EXPIRY")

	viper.BindEnv("rabbitmq.host", "RABBITMQ_HOST")
	viper.BindEnv("rabbitmq.port", "RABBITMQ_PORT")
	viper.BindEnv("rabbitmq.user", "RABBITMQ_USER")
	viper.BindEnv("rabbitmq.password", "RABBITMQ_PASSWORD")
	viper.BindEnv("rabbitmq.exchange", "RABBITMQ_EXCHANGE")

	viper.BindEnv("smtp.host", "SMTP_HOST")
	viper.BindEnv("smtp.port", "SMTP_PORT")
	viper.BindEnv("smtp.user", "SMTP_USER")
	viper.BindEnv("smtp.password", "SMTP_PASSWORD")
	viper.BindEnv("smtp.from", "SMTP_FROM")

	viper.BindEnv("oauth.google.web_client_id", "GOOGLE_WEB_CLIENT_ID")
	viper.BindEnv("oauth.google.android_client_id", "GOOGLE_ANDROID_CLIENT_ID")
	viper.BindEnv("oauth.google.ios_client_id", "GOOGLE_IOS_CLIENT_ID")
	viper.BindEnv("oauth.google.redirect_url", "GOOGLE_REDIRECT_URL")

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

	config.Database.Host = viper.GetString("database.host")
	config.Database.Port = viper.GetString("database.port")
	config.Database.User = viper.GetString("database.user")
	config.Database.Password = viper.GetString("database.password")
	config.Database.Name = viper.GetString("database.name")
	config.Database.SSLMode = viper.GetString("database.sslmode")

	config.JWT.Secret = viper.GetString("jwt.secret")
	expiry, err := time.ParseDuration(viper.GetString("jwt.expiry"))
	if err != nil {
		expiry = 24 * time.Hour
	}
	config.JWT.Expiry = expiry

	config.RabbitMQ.Host = viper.GetString("rabbitmq.host")
	config.RabbitMQ.Port = viper.GetString("rabbitmq.port")
	config.RabbitMQ.User = viper.GetString("rabbitmq.user")
	config.RabbitMQ.Password = viper.GetString("rabbitmq.password")
	config.RabbitMQ.Exchange = viper.GetString("rabbitmq.exchange")

	config.SMTP.Host = viper.GetString("smtp.host")
	config.SMTP.Port = viper.GetString("smtp.port")
	config.SMTP.User = viper.GetString("smtp.user")
	config.SMTP.Password = viper.GetString("smtp.password")
	config.SMTP.From = viper.GetString("smtp.from")

	config.GoogleAuth.WebClientID = viper.GetString("oauth.google.web_client_id")
	config.GoogleAuth.AndroidClientID = viper.GetString("oauth.google.android_client_id")
	config.GoogleAuth.IOSClientID = viper.GetString("oauth.google.ios_client_id")
	config.GoogleAuth.RedirectURL = viper.GetString("oauth.google.redirect_url")

	return &config, nil
}
