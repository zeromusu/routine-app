package persistence

import (
	"routine-app-server/internal/domain/repository"

	"gorm.io/gorm"
)

type routinePersistence struct {
	db *gorm.DB
}

func NewRoutinePersistence(db *gorm.DB) repository.RoutineRepository {
	return &routinePersistence{db: db}
}
