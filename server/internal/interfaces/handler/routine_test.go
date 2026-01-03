package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"routine-app-server/internal/domain"
	"routine-app-server/internal/interfaces/request"
	"routine-app-server/internal/interfaces/response"
	"routine-app-server/internal/usecase/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRoutineHandlerCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockUC := mocks.NewRoutineUseCase(t)
		h := NewRoutineHandler(mockUC)

		input := request.CreateRoutineRequest{
			Title:    "筋トレ",
			Interval: "daily",
		}
		body, _ := json.Marshal(input)

		expectedRoutine := &domain.Routine{Title: "筋トレ", Interval: "daily"}
		mockUC.On("CreateRoutine", input.Title, input.Interval).Return(expectedRoutine, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/v1/routines", bytes.NewBuffer(body))

		h.Create(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		var resp response.APIResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp.Success)

		data := resp.Data.(map[string]interface{})
		assert.Equal(t, expectedRoutine.Title, data["title"])
		assert.Equal(t, expectedRoutine.Interval, data["interval"])
	})

	t.Run("Invalid", func(t *testing.T) {
		mockUC := mocks.NewRoutineUseCase(t)
		h := NewRoutineHandler(mockUC)

		input := request.CreateRoutineRequest{
			Title:    "不正",
			Interval: "invalid",
		}
		body, _ := json.Marshal(input)

		mockUC.On("CreateRoutine", input.Title, input.Interval).Return(nil, domain.ErrInvalidData)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/v1/routines", bytes.NewBuffer(body))

		h.Create(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp response.APIResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, string(response.CodeInvalidPayload), resp.Error.Code)
	})
	t.Run("Duplicate", func(t *testing.T) {
		mockUC := mocks.NewRoutineUseCase(t)
		h := NewRoutineHandler(mockUC)

		input := request.CreateRoutineRequest{
			Title:    "重複テスト",
			Interval: "daily",
		}
		body, _ := json.Marshal(input)

		mockUC.On("CreateRoutine", "重複テスト", "daily").Return(nil, domain.ErrDuplicate)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/v1/routines", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Create(c)

		assert.Equal(t, http.StatusConflict, w.Code)

		var resp response.APIResponse
		json.Unmarshal(w.Body.Bytes(), &resp)

		assert.False(t, resp.Success)
		assert.Equal(t, string(response.CodeDuplicateRoutine), resp.Error.Code)
	})

	t.Run("Database Error", func(t *testing.T) {
		mockUC := mocks.NewRoutineUseCase(t)
		h := NewRoutineHandler(mockUC)

		input := request.CreateRoutineRequest{
			Title:    "DB Error",
			Interval: "db error",
		}
		body, _ := json.Marshal(input)

		mockUC.On("CreateRoutine", mock.Anything, mock.Anything).Return(nil, domain.ErrDatabase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/v1/routines", bytes.NewBuffer(body))

		h.Create(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var resp response.APIResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, string(response.CodeInternalServerError), resp.Error.Code)
	})
}
