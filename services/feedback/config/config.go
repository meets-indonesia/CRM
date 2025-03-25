package config

import (
	"github.com/spf13/viper"
)

// Config berisi semua konfigurasi aplikasi
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	RabbitMQ  RabbitMQConfig
	JWT       JWTConfig
	Email     EmailConfig
	FileStore FileStoreConfig
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

// RabbitMQConfig untuk konfigurasi RabbitMQ
type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Exchange string
}

// JWTConfig untuk konfigurasi JWT
type JWTConfig struct {
	Secret string
}

// EmailConfig untuk konfigurasi email
type EmailConfig struct {
	Host       string
	Port       string
	Username   string
	Password   string
	From       string
	AdminEmail string
}

// FileStoreConfig untuk konfigurasi penyimpanan file
type FileStoreConfig struct {
	UploadDir string
	MaxSize   int64 // dalam bytes
}

// LoadConfig memuat konfigurasi dari environment atau file
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	// Default values
	viper.SetDefault("server.port", "8083")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.host", "feedback-db")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.name", "feedback_db")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("rabbitmq.host", "rabbitmq")
	viper.SetDefault("rabbitmq.port", "5672")
	viper.SetDefault("rabbitmq.user", "guest")
	viper.SetDefault("rabbitmq.password", "guest")
	viper.SetDefault("rabbitmq.exchange", "feedback.events")
	viper.SetDefault("jwt.secret", "your-secret-key")

	// Default email config
	viper.SetDefault("email.host", "smtp.gmail.com")
	viper.SetDefault("email.port", "587")
	viper.SetDefault("email.username", "")
	viper.SetDefault("email.password", "")
	viper.SetDefault("email.from", "noreply@lrtsumsel.id")
	viper.SetDefault("email.admin_email", "adminlrt@gmail.com")

	// Default file store config
	viper.SetDefault("filestore.upload_dir", "./uploads")
	viper.SetDefault("filestore.max_size", 5*1024*1024) // 5MB

	// Environment variables mapping
	viper.BindEnv("server.port", "PORT")
	viper.BindEnv("server.mode", "GIN_MODE")
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.name", "DB_NAME")
	viper.BindEnv("database.sslmode", "DB_SSLMODE")
	viper.BindEnv("rabbitmq.host", "RABBITMQ_HOST")
	viper.BindEnv("rabbitmq.port", "RABBITMQ_PORT")
	viper.BindEnv("rabbitmq.user", "RABBITMQ_USER")
	viper.BindEnv("rabbitmq.password", "RABBITMQ_PASSWORD")
	viper.BindEnv("rabbitmq.exchange", "RABBITMQ_EXCHANGE")
	viper.BindEnv("jwt.secret", "JWT_SECRET")

	// Email environment variables
	viper.BindEnv("email.host", "EMAIL_HOST")
	viper.BindEnv("email.port", "EMAIL_PORT")
	viper.BindEnv("email.username", "EMAIL_USERNAME")
	viper.BindEnv("email.password", "EMAIL_PASSWORD")
	viper.BindEnv("email.from", "EMAIL_FROM")
	viper.BindEnv("email.admin_email", "ADMIN_EMAIL")

	// File store environment variables
	viper.BindEnv("filestore.upload_dir", "UPLOAD_DIR")
	viper.BindEnv("filestore.max_size", "MAX_UPLOAD_SIZE")

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

	config.RabbitMQ.Host = viper.GetString("rabbitmq.host")
	config.RabbitMQ.Port = viper.GetString("rabbitmq.port")
	config.RabbitMQ.User = viper.GetString("rabbitmq.user")
	config.RabbitMQ.Password = viper.GetString("rabbitmq.password")
	config.RabbitMQ.Exchange = viper.GetString("rabbitmq.exchange")

	config.JWT.Secret = viper.GetString("jwt.secret")

	// Set email config
	config.Email.Host = viper.GetString("email.host")
	config.Email.Port = viper.GetString("email.port")
	config.Email.Username = viper.GetString("email.username")
	config.Email.Password = viper.GetString("email.password")
	config.Email.From = viper.GetString("email.from")
	config.Email.AdminEmail = viper.GetString("email.admin_email")

	// Set file store config
	config.FileStore.UploadDir = viper.GetString("filestore.upload_dir")
	config.FileStore.MaxSize = viper.GetInt64("filestore.max_size")

	return &config, nil
}
