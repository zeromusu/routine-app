package persistence

import (
	"errors"
	"regexp"
	"routine-app-server/internal/domain"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestFindAll(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	repo := NewRoutinePersistence(gormDB)

	getAllSQL := regexp.QuoteMeta(`SELECT * FROM "routines"`)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "interval"}).
			AddRow(1, "筋トレ", "daily").
			AddRow(2, "読書", "weekly")

		mock.ExpectQuery(getAllSQL).WillReturnRows(rows)
		results, err := repo.FindAll()

		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, "筋トレ", results[0].Title)
	})

	t.Run("Database Error", func(t *testing.T) {
		mock.ExpectQuery(getAllSQL).WillReturnError(errors.New("db error"))

		results, err := repo.FindAll()

		assert.Nil(t, results)
		assert.ErrorIs(t, err, domain.ErrDatabase)
	})
}

func TestCreate(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	repo := NewRoutinePersistence(gormDB)

	insertSQL := regexp.QuoteMeta(`INSERT INTO "routines`)

	t.Run("Success", func(t *testing.T) {
		routine := &domain.Routine{Title: "Reading", Interval: "daily"}

		mock.ExpectBegin()

		mock.ExpectQuery(insertSQL).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		err := repo.Create(routine)

		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Invalid Data", func(t *testing.T) {
		routine := &domain.Routine{Title: ""}

		mock.ExpectBegin()
		mock.ExpectQuery(insertSQL).WillReturnError(gorm.ErrInvalidData)
		mock.ExpectRollback()

		err := repo.Create(routine)

		assert.ErrorIs(t, err, domain.ErrInvalidData)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("23505 -> Duplicate", func(t *testing.T) {
		routine := &domain.Routine{Title: "重複テスト"}

		pgErr := &pgconn.PgError{Code: "23505"}

		mock.ExpectBegin()

		mock.ExpectQuery(insertSQL).WillReturnError(pgErr)
		mock.ExpectRollback()

		err := repo.Create(routine)

		assert.True(t, errors.Is(err, domain.ErrDuplicate))
	})

	t.Run("Database Error", func(t *testing.T) {
		routine := &domain.Routine{Title: "Invalid"}

		mock.ExpectBegin()
		mock.ExpectQuery(insertSQL).WillReturnError(errors.New("connection reset by peer"))
		mock.ExpectRollback()

		err := repo.Create(routine)

		assert.ErrorIs(t, err, domain.ErrDatabase)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
