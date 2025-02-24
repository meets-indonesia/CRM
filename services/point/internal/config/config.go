package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	// Database configurations
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`

	// Server configurations
	ServerPort string `mapstructure:"SERVER_PORT"`

	// JWT configuration
	JWTSecret string `mapstructure:"JWT_SECRET"`

	// RabbitMQ configuration
	RabbitMQURL string `mapstructure:"RABBITMQ_URL"`

	// Auth Service configuration
	AuthServiceURL string `mapstructure:"AUTH_SERVICE_URL"`
}

func LoadConfig() (Config, error) {
	var config Config

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "point_db")
	viper.SetDefault("SERVER_PORT", "8083")
	viper.SetDefault("JWT_SECRET", "your-secret-key")
	viper.SetDefault("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/")
	viper.SetDefault("AUTH_SERVICE_URL", "http://auth:8080")

	viper.AutomaticEnv()

	config.DBHost = viper.GetString("DB_HOST")
	config.DBPort = viper.GetString("DB_PORT")
	config.DBUser = viper.GetString("DB_USER")
	config.DBPassword = viper.GetString("DB_PASSWORD")
	config.DBName = viper.GetString("DB_NAME")
	config.ServerPort = viper.GetString("SERVER_PORT")
	config.JWTSecret = viper.GetString("JWT_SECRET")
	config.RabbitMQURL = viper.GetString("RABBITMQ_URL")
	config.AuthServiceURL = viper.GetString("AUTH_SERVICE_URL")

	return config, nil
}
