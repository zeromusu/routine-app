package usecase

import (
	"errors"
	"routine-app-server/internal/domain"
	"routine-app-server/internal/domain/repository/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRoutineUseCaseGetRoutines(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := mocks.NewRoutineRepository(t)
		uc := NewRoutineUseCase(mockRepo)

		expected := []*domain.Routine{{Title: "Test Routine", Interval: "daily"}}
		mockRepo.On("FindAll").Return(expected, nil)

		res, err := uc.GetRoutines()

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
		assert.Len(t, res, 1)
	})

	t.Run("Database Error", func(t *testing.T) {
		mockRepo := mocks.NewRoutineRepository(t)
		uc := NewRoutineUseCase(mockRepo)

		mockRepo.On("FindAll").Return(nil, errors.New("db error"))

		res, err := uc.GetRoutines()

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestRoutineUseCaseGetRoutine(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := mocks.NewRoutineRepository(t)
		uc := NewRoutineUseCase(mockRepo)

		id := 1

		expected := &domain.Routine{ID: id, Title: "Test Routine", Interval: "daily"}
		mockRepo.On("FindOne", id).Return(expected, nil)

		res, err := uc.GetRoutine(id)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo := mocks.NewRoutineRepository(t)
		uc := NewRoutineUseCase(mockRepo)

		id := 1
		mockRepo.On("FindOne", id).Return(nil, domain.ErrNotFound)

		res, err := uc.GetRoutine(id)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("Database Error", func(t *testing.T) {
		mockRepo := mocks.NewRoutineRepository(t)
		uc := NewRoutineUseCase(mockRepo)

		id := 1
		mockRepo.On("FindOne", id).Return(nil, domain.ErrDatabase)

		res, err := uc.GetRoutine(id)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.ErrorIs(t, err, domain.ErrDatabase)
	})
}

func TestRoutineUseCaseCreateRoutine(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := mocks.NewRoutineRepository(t)
		uc := NewRoutineUseCase(mockRepo)

		title := "自動生成テスト"
		interval := "weekly"

		mockRepo.On("Create", mock.MatchedBy(func(r *domain.Routine) bool {
			return r.Title == title && r.Interval == interval
		})).Return(nil)

		res, err := uc.CreateRoutine(title, interval)

		assert.NoError(t, err)
		assert.Equal(t, title, res.Title)
		assert.Equal(t, interval, res.Interval)
	})

	t.Run("Invalid Data", func(t *testing.T) {
		mockRepo := mocks.NewRoutineRepository(t)
		uc := NewRoutineUseCase(mockRepo)

		mockRepo.On("Create", mock.Anything).Return(domain.ErrInvalidData)

		res, err := uc.CreateRoutine("不正", "invalid")

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.ErrorIs(t, err, domain.ErrInvalidData)
	})

	t.Run("Duplicate Data", func(t *testing.T) {
		mockRepo := mocks.NewRoutineRepository(t)
		uc := NewRoutineUseCase(mockRepo)

		mockRepo.On("Create", mock.Anything).Return(domain.ErrDuplicate)

		res, err := uc.CreateRoutine("重複タイトル", "daily")

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.ErrorIs(t, err, domain.ErrDuplicate)
	})

	t.Run("Database Error", func(t *testing.T) {
		mockRepo := mocks.NewRoutineRepository(t)
		uc := NewRoutineUseCase(mockRepo)

		mockRepo.On("Create", mock.Anything).Return(domain.ErrDatabase)

		res, err := uc.CreateRoutine("test", "daily")

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.ErrorIs(t, err, domain.ErrDatabase)
	})
}
