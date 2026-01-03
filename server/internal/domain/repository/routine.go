package repository

import "routine-app-server/internal/domain"

//go:generate mockery --name=RoutineRepository --case=underscore --output=./mocks

type RoutineRepository interface {
	Create(routine *domain.Routine) error
}
