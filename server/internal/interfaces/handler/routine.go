package handler

import (
	"errors"
	"net/http"
	"routine-app-server/internal/domain"
	"routine-app-server/internal/interfaces/request"
	"routine-app-server/internal/interfaces/response"
	"routine-app-server/internal/usecase"

	"github.com/gin-gonic/gin"
)

type RoutineHandler interface {
	GetAll(c *gin.Context)
	Create(c *gin.Context)
}

type routineHandler struct {
	routineUseCase usecase.RoutineUseCase
}

func NewRoutineHandler(uc usecase.RoutineUseCase) RoutineHandler {
	return &routineHandler{
		routineUseCase: uc,
	}
}

func (h *routineHandler) GetAll(c *gin.Context) {
	routines, err := h.routineUseCase.GetRoutines()
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, string(response.CodeInternalServerError), err.Error())
		return
	}

	response.RespondSuccess(c, http.StatusOK, routines)
}

func (h *routineHandler) Create(c *gin.Context) {
	var createRoutine request.CreateRoutineRequest
	if err := c.ShouldBindJSON(&createRoutine); err != nil {
		response.RespondError(c, http.StatusBadRequest, string(response.CodeInvalidPayload), err.Error())
		return
	}

	routine, err := h.routineUseCase.CreateRoutine(createRoutine.Title, createRoutine.Interval)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicate) {
			response.RespondError(c, http.StatusConflict, string(response.CodeDuplicateRoutine), err.Error())
			return
		}
		if errors.Is(err, domain.ErrInvalidData) {
			response.RespondError(c, http.StatusBadRequest, string(response.CodeInvalidPayload), err.Error())
			return
		}
		response.RespondError(c, http.StatusInternalServerError, string(response.CodeInternalServerError), err.Error())
		return
	}

	response.RespondSuccess(c, http.StatusCreated, routine)
}
