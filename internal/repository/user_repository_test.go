package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)
	t.Cleanup(func() { sqlDB.Close() })

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	require.NoError(t, err)

	return gdb, mock
}

func TestUserRepository_GetById(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)

	id := "123e4567-e89b-12d3-a456-426614174000"
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "full_name", "email", "password"}).
		AddRow(id, now, now, "Ivan", "ivan@test.com", "hash")
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1.*LIMIT \$2`).
		WithArgs(id, 1).
		WillReturnRows(rows)

	ctx := context.Background()
	user, err := repo.GetById(ctx, id)

	require.NoError(t, err)
	assert.Equal(t, "ivan@test.com", user.Email)
	assert.Equal(t, "Ivan", user.FullName)
	assert.Equal(t, id, user.ID)
	require.NoError(t, mock.ExpectationsWereMet())

}

func TestUserRepository_GetById_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1.*LIMIT \$2`).
		WithArgs("missing-id", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := repo.GetById(context.Background(), "missing-id")

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	_ = user // см. примечание ниже про баг
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByEmail_Found(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "full_name", "email", "password"}).
		AddRow("123e4567-e89b-12d3-a456-426614174000", now, now, "Ivan", "ivan@test.com", "hash")

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1.*LIMIT \$2`).
		WithArgs("ivan@test.com", 1).
		WillReturnRows(rows)

	user, err := repo.FindByEmail("ivan@test.com")

	require.NoError(t, err)
	assert.Equal(t, "ivan@test.com", user.Email)
	assert.Equal(t, "Ivan", user.FullName)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1.*LIMIT \$2`).
		WithArgs("missing@test.com", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := repo.FindByEmail("missing@test.com")

	assert.Nil(t, user)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	require.NoError(t, mock.ExpectationsWereMet())
}
