package usecase

import "routine-app-server/internal/domain/repository"

type RoutineUseCase interface {
}

type routineUseCase struct {
	routineRepository repository.RoutineRepository
}

func NewRoutineUseCase(repo repository.RoutineRepository) RoutineUseCase {
	return &routineUseCase{
		routineRepository: repo,
	}
}
