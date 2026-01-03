package usecase

import (
	"routine-app-server/internal/domain"
	"routine-app-server/internal/domain/repository"
)

//go:generate mockery --all --case=underscore --outpkg=mocks --output=./mocks

type RoutineUseCase interface {
	CreateRoutine(title, interval string) (*domain.Routine, error)
}

type routineUseCase struct {
	routineRepository repository.RoutineRepository
}

func NewRoutineUseCase(repo repository.RoutineRepository) RoutineUseCase {
	return &routineUseCase{
		routineRepository: repo,
	}
}

func (u *routineUseCase) CreateRoutine(title, interval string) (*domain.Routine, error) {
	routine := &domain.Routine{
		Title:    title,
		Interval: interval,
	}

	if err := u.routineRepository.Create(routine); err != nil {
		return nil, err
	}
	return routine, nil
}
