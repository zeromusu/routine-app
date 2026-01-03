package handler

import "routine-app-server/internal/usecase"

type RoutineHandler interface {
}

type routineHandler struct {
	routineUseCase usecase.RoutineUseCase
}

func NewRoutineHandler(uc usecase.RoutineUseCase) RoutineHandler {
	return &routineHandler{
		routineUseCase: uc,
	}
}
