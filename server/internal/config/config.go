package config

import "os"

type Config struct {
	DB DBConfig
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
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Port:     os.Getenv("DB_PORT"),
		},
	}
}
