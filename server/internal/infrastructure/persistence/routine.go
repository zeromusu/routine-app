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

func (p *routinePersistence) FindAll() ([]*domain.Routine, error) {
	var routines []*domain.Routine

	if err := p.db.Find(&routines).Error; err != nil {
		return nil, domain.ErrDatabase
	}
	return routines, nil
}

func (p *routinePersistence) FindOne(id int) (*domain.Routine, error) {
	var routine domain.Routine

	if err := p.db.First(&routine, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, domain.ErrDatabase
	}
	return &routine, nil
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
