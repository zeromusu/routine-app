package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestRoutineHandlerGetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockUC := mocks.NewRoutineUseCase(t)
		h := NewRoutineHandler(mockUC)

		expected := []*domain.Routine{
			{ID: 1, Title: "Routine 1", Interval: "daily"},
			{ID: 2, Title: "Routine 2", Interval: "weekly"},
		}
		mockUC.On("GetRoutines").Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		h.GetAll(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Routine 1")
		assert.Contains(t, w.Body.String(), "Routine 2")
	})

	t.Run("Database Error", func(t *testing.T) {
		mockUC := mocks.NewRoutineUseCase(t)
		h := NewRoutineHandler(mockUC)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		mockUC.On("GetRoutines").Return(nil, domain.ErrDatabase)

		h.GetAll(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestRoutineHandlerGetOne(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockUC := mocks.NewRoutineUseCase(t)
		h := NewRoutineHandler(mockUC)

		id := 1

		expected := &domain.Routine{
			ID: id, Title: "Routine 1", Interval: "daily",
		}
		mockUC.On("GetRoutine", id).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest("GET", fmt.Sprintf("/v1/routines/%d", id), nil)

		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", id)}}

		h.GetOne(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp response.APIResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp.Success)

		data := resp.Data.(map[string]interface{})
		assert.Equal(t, float64(id), data["id"])
		assert.Equal(t, "Routine 1", data["title"])
	})

	t.Run("Not Found", func(t *testing.T) {
		mockUC := mocks.NewRoutineUseCase(t)
		h := NewRoutineHandler(mockUC)

		id := 99
		mockUC.On("GetRoutine", id).Return(nil, domain.ErrNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", fmt.Sprintf("/v1/routines/%d", id), nil)
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", id)}}

		h.GetOne(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Database Error", func(t *testing.T) {
		mockUC := mocks.NewRoutineUseCase(t)
		h := NewRoutineHandler(mockUC)

		id := 1
		mockUC.On("GetRoutine", id).Return(nil, domain.ErrDatabase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", fmt.Sprintf("/v1/routines/%d", id), nil)
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", id)}}

		h.GetOne(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestRoutineHandlerCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
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

	t.Run("Invalid Data", func(t *testing.T) {
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
	t.Run("Duplicate Data", func(t *testing.T) {
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

func TestRoutineHandlerUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockUC := mocks.NewRoutineUseCase(t)
		h := NewRoutineHandler(mockUC)

		input := request.UpdateRoutineRequest{
			Title:    "更新後タイトル",
			Interval: "daily",
		}
		body, _ := json.Marshal(input)

		id := 1

		expectedRoutine := &domain.Routine{ID: id, Title: "更新後タイトル", Interval: "daily"}
		mockUC.On("UpdateRoutine", id, input.Title, input.Interval).Return(expectedRoutine, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", fmt.Sprintf("/v1/routines/%d", id), bytes.NewBuffer(body))
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", id)}}

		h.Update(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp response.APIResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp.Success)

		data := resp.Data.(map[string]interface{})
		assert.Equal(t, expectedRoutine.Title, data["title"])
		assert.Equal(t, expectedRoutine.Interval, data["interval"])
	})

	t.Run("Invalid Data", func(t *testing.T) {
		mockUC := mocks.NewRoutineUseCase(t)
		h := NewRoutineHandler(mockUC)

		input := request.UpdateRoutineRequest{
			Title:    "不正",
			Interval: "invalid",
		}
		body, _ := json.Marshal(input)

		id := 1

		mockUC.On("UpdateRoutine", id, input.Title, input.Interval).Return(nil, domain.ErrInvalidData)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", fmt.Sprintf("/v1/routines/%d", id), bytes.NewBuffer(body))
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", id)}}

		h.Update(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp response.APIResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, string(response.CodeInvalidPayload), resp.Error.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockUC := mocks.NewRoutineUseCase(t)
		h := NewRoutineHandler(mockUC)

		input := request.UpdateRoutineRequest{
			Title:    "test",
			Interval: "daily",
		}
		body, _ := json.Marshal(input)

		id := 99
		mockUC.On("UpdateRoutine", id, input.Title, input.Interval).Return(nil, domain.ErrNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", fmt.Sprintf("/v1/routines/%d", id), bytes.NewBuffer(body))
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", id)}}

		h.Update(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var resp response.APIResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, string(response.CodeNotFound), resp.Error.Code)
	})

	t.Run("Database Error", func(t *testing.T) {
		mockUC := mocks.NewRoutineUseCase(t)
		h := NewRoutineHandler(mockUC)

		input := request.UpdateRoutineRequest{
			Title:    "test",
			Interval: "daily",
		}
		body, _ := json.Marshal(input)

		id := 1
		mockUC.On("UpdateRoutine", id, input.Title, input.Interval).Return(nil, domain.ErrDatabase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", fmt.Sprintf("/v1/routines/%d", id), bytes.NewBuffer(body))
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", id)}}

		h.Update(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var resp response.APIResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, string(response.CodeInternalServerError), resp.Error.Code)
	})
}

func TestRoutineHandlerDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockUC := mocks.NewRoutineUseCase(t)
		h := NewRoutineHandler(mockUC)

		id := 1

		mockUC.On("DeleteRoutine", id).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("DELETE", fmt.Sprintf("/v1/routines/%d", id), nil)
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", id)}}

		h.Delete(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp response.APIResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp.Success)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockUC := mocks.NewRoutineUseCase(t)
		h := NewRoutineHandler(mockUC)

		id := 99

		mockUC.On("DeleteRoutine", id).Return(domain.ErrNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("DELETE", fmt.Sprintf("/v1/routines/%d", id), nil)
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", id)}}

		h.Delete(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var resp response.APIResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, string(response.CodeNotFound), resp.Error.Code)
	})

	t.Run("Database Error", func(t *testing.T) {
		mockUC := mocks.NewRoutineUseCase(t)
		h := NewRoutineHandler(mockUC)

		id := 1

		mockUC.On("DeleteRoutine", id).Return(domain.ErrDatabase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("DELETE", fmt.Sprintf("/v1/routines/%d", id), nil)
		c.Params = gin.Params{{Key: "id", Value: fmt.Sprintf("%d", id)}}

		h.Delete(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var resp response.APIResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, string(response.CodeInternalServerError), resp.Error.Code)
	})
}
