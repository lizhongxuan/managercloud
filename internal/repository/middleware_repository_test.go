package repository

import (
	"testing"
	"time"

	"middleware-platform/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm DB: %v", err)
	}

	return db, mock
}

func TestMiddlewareRepository_Create(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewMiddlewareRepository(db)

	mw := &model.Middleware{
		Name:    "test-redis",
		Type:    "redis",
		Version: "6.2",
		Host:    "localhost",
		Port:    "6379",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "middlewares"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(mw)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMiddlewareRepository_FindByType(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewMiddlewareRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "type", "version", "host", "port", "created_at", "updated_at"}).
		AddRow(1, "test-redis", "redis", "6.2", "localhost", "6379", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "middlewares"`).
		WithArgs("redis").
		WillReturnRows(rows)

	middlewares, err := repo.FindByType("redis")
	assert.NoError(t, err)
	assert.Len(t, middlewares, 1)
	assert.Equal(t, "test-redis", middlewares[0].Name)
} 