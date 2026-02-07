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

func TestRoutinePersistenceFindAll(t *testing.T) {
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

func TestRoutinePersistenceFindOne(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	repo := NewRoutinePersistence(gormDB)

	getOneSQL := regexp.QuoteMeta(`SELECT * FROM "routines" WHERE "routines"."id" = $1 ORDER BY "routines"."id" LIMIT $2`)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "interval"}).
			AddRow(1, "筋トレ", "daily")

		mock.ExpectQuery(getOneSQL).
			WithArgs(1, 1).
			WillReturnRows(rows)
		result, err := repo.FindOne(1)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.ID)
		assert.Equal(t, "筋トレ", result.Title)
		assert.Equal(t, "daily", result.Interval)
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery(getOneSQL).
			WithArgs(99, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "interval"}))

		result, err := repo.FindOne(99)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("Database Error", func(t *testing.T) {
		mock.ExpectQuery(getOneSQL).WillReturnError(errors.New("db error"))

		results, err := repo.FindAll()

		assert.Nil(t, results)
		assert.ErrorIs(t, err, domain.ErrDatabase)
	})
}

func TestRoutinePersistenceCreate(t *testing.T) {
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

func TestRoutinePersistenceUpdate(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	repo := NewRoutinePersistence(gormDB)

	updateSQL := `UPDATE "routines" SET`

	t.Run("Success", func(t *testing.T) {
		routine := &domain.Routine{
			ID:       1,
			Title:    "更新後",
			Interval: "weekly",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(updateSQL)).
			WithArgs(routine.Title, routine.Interval, sqlmock.AnyArg(), routine.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(routine)

		assert.NoError(t, err)
	})

	t.Run("Invalid Data", func(t *testing.T) {
		routine := &domain.Routine{
			ID:       1,
			Title:    "invalid",
			Interval: "invalid",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(updateSQL)).
			WithArgs(routine.Title, routine.Interval, sqlmock.AnyArg(), routine.ID).
			WillReturnError(gorm.ErrInvalidData)
		mock.ExpectRollback()

		err := repo.Update(routine)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidData)
	})

	t.Run("Not Found", func(t *testing.T) {
		routine := &domain.Routine{ID: 99}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(updateSQL)).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := repo.Update(routine)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("Database Error", func(t *testing.T) {
		routine := &domain.Routine{
			ID:       1,
			Title:    "更新後",
			Interval: "weekly",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(updateSQL)).
			WillReturnError(errors.New("connection reset by peer"))
		mock.ExpectRollback()

		err := repo.Update(routine)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrDatabase)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRoutinePersistenceDelete(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	repo := NewRoutinePersistence(gormDB)

	deleteSQL := regexp.QuoteMeta(`DELETE FROM "routines" WHERE "routines"."id" = $1`)

	t.Run("Success", func(t *testing.T) {
		id := 1
		mock.ExpectBegin()
		mock.ExpectExec(deleteSQL).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.Delete(id)
		assert.NoError(t, err)
	})

	t.Run("Not Found", func(t *testing.T) {
		id := 99
		mock.ExpectBegin()
		mock.ExpectExec(deleteSQL).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		err := repo.Delete(id)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("Database Error", func(t *testing.T) {
		id := 1
		mock.ExpectBegin()
		mock.ExpectExec(deleteSQL).
			WithArgs(id).
			WillReturnError(errors.New("connection reset by peer"))
		mock.ExpectRollback()

		err := repo.Delete(id)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrDatabase)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
