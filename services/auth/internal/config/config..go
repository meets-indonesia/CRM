package config

import (
	"time"

	"github.com/spf13/viper"
)

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

	return config, nil
}
