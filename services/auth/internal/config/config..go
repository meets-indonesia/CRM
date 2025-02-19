// services/auth/internal/config/config.go
package config

import (
	"time"

	"github.com/spf13/viper"
)

type OAuthConfig struct {
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	AppleClientID      string `mapstructure:"APPLE_CLIENT_ID"`
	AppleClientSecret  string `mapstructure:"APPLE_CLIENT_SECRET"`
	AppleKeyID         string `mapstructure:"APPLE_KEY_ID"`
	AppleTeamID        string `mapstructure:"APPLE_TEAM_ID"`
	ApplePrivateKey    string `mapstructure:"APPLE_PRIVATE_KEY"`
}

type Config struct {
	DBHost        string `mapstructure:"DB_HOST"`
	DBPort        string `mapstructure:"DB_PORT"`
	DBUser        string `mapstructure:"DB_USER"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBName        string `mapstructure:"DB_NAME"`
	JWTSecret     string `mapstructure:"JWT_SECRET"`
	JWTExpiration time.Duration
	ServerPort    string `mapstructure:"SERVER_PORT"`
	SMTPHost      string `mapstructure:"SMTP_HOST"`
	SMTPPort      int    `mapstructure:"SMTP_PORT"`
	SMTPUsername  string `mapstructure:"SMTP_USERNAME"`
	SMTPPassword  string `mapstructure:"SMTP_PASSWORD"`
	RabbitMQURL   string `mapstructure:"RABBITMQ_URL"`
	// Add OAuth configuration
	OAuth OAuthConfig
}

func LoadConfig() (Config, error) {
	var config Config

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "auth_db")
	viper.SetDefault("JWT_SECRET", "your-secret-key")
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SMTP_HOST", "smtp.gmail.com")
	viper.SetDefault("SMTP_PORT", 587)
	viper.SetDefault("SMTP_USERNAME", "your-email@gmail.com")
	viper.SetDefault("SMTP_PASSWORD", "your-app-password")
	viper.SetDefault("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/")
	viper.SetDefault("GOOGLE_CLIENT_ID", "")
	viper.SetDefault("GOOGLE_CLIENT_SECRET", "")
	viper.SetDefault("APPLE_CLIENT_ID", "")
	viper.SetDefault("APPLE_CLIENT_SECRET", "")
	viper.SetDefault("APPLE_KEY_ID", "")
	viper.SetDefault("APPLE_TEAM_ID", "")
	viper.SetDefault("APPLE_PRIVATE_KEY", "")

	viper.AutomaticEnv()

	config.DBHost = viper.GetString("DB_HOST")
	config.DBPort = viper.GetString("DB_PORT")
	config.DBUser = viper.GetString("DB_USER")
	config.DBPassword = viper.GetString("DB_PASSWORD")
	config.DBName = viper.GetString("DB_NAME")
	config.JWTSecret = viper.GetString("JWT_SECRET")
	config.JWTExpiration = 24 * time.Hour
	config.ServerPort = viper.GetString("SERVER_PORT")
	config.SMTPHost = viper.GetString("SMTP_HOST")
	config.SMTPPort = viper.GetInt("SMTP_PORT")
	config.SMTPUsername = viper.GetString("SMTP_USERNAME")
	config.SMTPPassword = viper.GetString("SMTP_PASSWORD")
	config.RabbitMQURL = viper.GetString("RABBITMQ_URL")
	config.OAuth.GoogleClientID = viper.GetString("GOOGLE_CLIENT_ID")
	config.OAuth.GoogleClientSecret = viper.GetString("GOOGLE_CLIENT_SECRET")
	config.OAuth.AppleClientID = viper.GetString("APPLE_CLIENT_ID")
	config.OAuth.AppleClientSecret = viper.GetString("APPLE_CLIENT_SECRET")
	config.OAuth.AppleKeyID = viper.GetString("APPLE_KEY_ID")
	config.OAuth.AppleTeamID = viper.GetString("APPLE_TEAM_ID")
	config.OAuth.ApplePrivateKey = viper.GetString("APPLE_PRIVATE_KEY")

	return config, nil
}
