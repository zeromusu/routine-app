package handler

import (
	"errors"
	"net/http"
	"routine-app-server/internal/domain"
	"routine-app-server/internal/interfaces/request"
	"routine-app-server/internal/interfaces/response"
	"routine-app-server/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoutineHandler interface {
	GetAll(c *gin.Context)
	GetOne(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type routineHandler struct {
	routineUseCase usecase.RoutineUseCase
}

func NewRoutineHandler(uc usecase.RoutineUseCase) RoutineHandler {
	return &routineHandler{
		routineUseCase: uc,
	}
}

// Create godoc
// @Summary Get all routines
// @Description Get all routines
// @Tags routines
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=[]domain.Routine}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /routines [get]
func (h *routineHandler) GetAll(c *gin.Context) {
	routines, err := h.routineUseCase.GetRoutines()
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, string(response.CodeInternalServerError), err.Error())
		return
	}

	response.RespondSuccess(c, http.StatusOK, routines)
}

// Create godoc
// @Summary Get one routine
// @Description Get one routine by ID
// @Tags routines
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=domain.Routine}
// @Failure 400 {object} response.APIResponse{error=response.APIError}
// @Failure 404 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /routines/{id} [get]
func (h *routineHandler) GetOne(c *gin.Context) {
	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		response.RespondError(c, http.StatusBadRequest, string(response.CodeInvalidPayload), err.Error())
		return
	}

	routine, err := h.routineUseCase.GetRoutine(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.RespondError(c, http.StatusNotFound, string(response.CodeNotFound), err.Error())
			return
		}
		response.RespondError(c, http.StatusInternalServerError, string(response.CodeInternalServerError), err.Error())
		return
	}

	response.RespondSuccess(c, http.StatusOK, routine)
}

// Create godoc
// @Summary Create a new routine
// @Description Create a new routine with the provided title and interval
// @Tags routines
// @Accept json
// @Produce json
// @Param routine body request.CreateRoutineRequest true "Routine to create"
// @Success 201 {object} response.APIResponse{data=domain.Routine}
// @Failure 400 {object} response.APIResponse{error=response.APIError}
// @Failure 409 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /routines [post]
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

// Create godoc
// @Summary Update a routine
// @Description Update a routine by ID
// @Tags routines
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=domain.Routine}
// @Failure 400 {object} response.APIResponse{error=response.APIError}
// @Failure 404 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /routines/{id} [put]
func (h *routineHandler) Update(c *gin.Context) {
	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		response.RespondError(c, http.StatusBadRequest, string(response.CodeInvalidPayload), err.Error())
		return
	}

	var updateRoutine request.UpdateRoutineRequest
	if err := c.ShouldBindJSON(&updateRoutine); err != nil {
		response.RespondError(c, http.StatusBadRequest, string(response.CodeInvalidPayload), err.Error())
		return
	}

	routine, err := h.routineUseCase.UpdateRoutine(id, updateRoutine.Title, updateRoutine.Interval)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidData) {
			response.RespondError(c, http.StatusBadRequest, string(response.CodeInvalidPayload), err.Error())
			return
		}
		if errors.Is(err, domain.ErrNotFound) {
			response.RespondError(c, http.StatusNotFound, string(response.CodeNotFound), err.Error())
			return
		}
		response.RespondError(c, http.StatusInternalServerError, string(response.CodeInternalServerError), err.Error())
		return
	}

	response.RespondSuccess(c, http.StatusOK, routine)
}

// Create godoc
// @Summary Delete a routine
// @Description Delete a routine by ID
// @Tags routines
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse{error=response.APIError}
// @Failure 404 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /routines/{id} [delete]
func (h *routineHandler) Delete(c *gin.Context) {
	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		response.RespondError(c, http.StatusBadRequest, string(response.CodeInvalidPayload), err.Error())
		return
	}

	if err := h.routineUseCase.DeleteRoutine(id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.RespondError(c, http.StatusNotFound, string(response.CodeNotFound), err.Error())
			return
		}
		response.RespondError(c, http.StatusInternalServerError, string(response.CodeInternalServerError), err.Error())
		return
	}

	response.RespondSuccess(c, http.StatusOK, gin.H{"message": "Routine deleted successfully"})
}
