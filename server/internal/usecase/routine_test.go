package usecase

import (
	"errors"
	"routine-app-server/internal/domain"
	"routine-app-server/internal/domain/repository/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRoutineUseCaseCreateRoutine(t *testing.T) {
	t.Run("success", func(t *testing.T) {
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

	t.Run("fail", func(t *testing.T) {
		mockRepo := mocks.NewRoutineRepository(t)
		uc := NewRoutineUseCase(mockRepo)

		mockRepo.On("Create", mock.Anything).Return(errors.New("db error"))

		res, err := uc.CreateRoutine("test", "daily")

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "db error", err.Error())
	})
}
