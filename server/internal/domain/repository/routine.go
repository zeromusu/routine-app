package repository

import "routine-app-server/internal/domain"

//go:generate mockery --name=RoutineRepository --case=underscore --output=./mocks

type RoutineRepository interface {
	FindAll() ([]*domain.Routine, error)
	FindOne(id int) (*domain.Routine, error)
	Create(routine *domain.Routine) error
}
