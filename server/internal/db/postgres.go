package db

import (
	"fmt"
	"routine-app-server/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(DBConfig config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		DBConfig.Host, DBConfig.User, DBConfig.Password, DBConfig.Name, DBConfig.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
