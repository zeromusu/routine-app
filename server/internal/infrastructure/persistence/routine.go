package persistence

import (
	"errors"
	"routine-app-server/internal/domain"
	"routine-app-server/internal/domain/repository"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type routinePersistence struct {
	db *gorm.DB
}

func NewRoutinePersistence(db *gorm.DB) repository.RoutineRepository {
	return &routinePersistence{db: db}
}

func (p *routinePersistence) Create(routine *domain.Routine) error {
	err := p.db.Create(routine).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrDuplicate
		}

		if errors.Is(err, gorm.ErrInvalidData) {
			return domain.ErrInvalidData
		}

		return domain.ErrDatabase
	}
	return nil
}
