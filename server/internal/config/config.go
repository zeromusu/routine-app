package config

import "os"

type Config struct {
	App AppConfig
	DB  DBConfig
}

type AppConfig struct {
	AppPort string
}

type DBConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
}

func LoadConfig() *Config {
	return &Config{
		App: AppConfig{
			AppPort: os.Getenv("APP_PORT"),
		},
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Port:     os.Getenv("DB_PORT"),
		},
	}
}
