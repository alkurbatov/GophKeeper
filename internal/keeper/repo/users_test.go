package repo_test

import (
	"context"
	"testing"

	"github.com/alkurbatov/goph-keeper/internal/keeper/entity"
	"github.com/alkurbatov/goph-keeper/internal/keeper/infra/postgres"
	"github.com/alkurbatov/goph-keeper/internal/keeper/repo"
	"github.com/alkurbatov/goph-keeper/internal/libraries/gophtest"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func newUsersRepo(t *testing.T, m pgxmock.PgxPoolIface) repo.Users {
	t.Helper()

	pg := &postgres.Postgres{
		Pool: m,
	}

	return repo.NewUsersRepo(pg, createTestLogger(t))
}

func TestRegisterUser(t *testing.T) {
	expected := uuid.NewV4()

	rows := pgxmock.NewRows([]string{"id"}).
		AddRow(expected.String())

	m := newPoolMock(t)
	m.ExpectBeginTx(postgres.DefaultTxOptions)
	m.ExpectQuery("INSERT INTO users").
		WithArgs(gophtest.Username, gophtest.SecurityKey).
		WillReturnRows(rows)
	m.ExpectCommit()

	sat := newUsersRepo(t, m)
	id, err := sat.Register(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.NoError(t, err)
	require.Equal(t, expected, id)
	require.NoError(t, m.ExpectationsWereMet())
}

func TestRegisterUserFailsOnUserExists(t *testing.T) {
	pErr := &pgconn.PgError{Code: pgerrcode.UniqueViolation}

	m := newPoolMock(t)
	m.ExpectBeginTx(postgres.DefaultTxOptions)
	m.ExpectQuery("INSERT INTO users").
		WithArgs(gophtest.Username, gophtest.SecurityKey).
		WillReturnError(error(pErr))
	m.ExpectRollback()

	sat := newUsersRepo(t, m)
	_, err := sat.Register(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.ErrorIs(t, err, entity.ErrUserExists)
	require.NoError(t, m.ExpectationsWereMet())
}

func TestRegisterUserFailsOnUnknownError(t *testing.T) {
	m := newPoolMock(t)
	m.ExpectBeginTx(postgres.DefaultTxOptions)
	m.ExpectQuery("INSERT INTO users").
		WithArgs(gophtest.Username, gophtest.SecurityKey).
		WillReturnError(gophtest.ErrUnexpected)
	m.ExpectRollback()

	sat := newUsersRepo(t, m)
	_, err := sat.Register(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.ErrorIs(t, err, gophtest.ErrUnexpected)
	require.NoError(t, m.ExpectationsWereMet())
}

func TestVerifyUser(t *testing.T) {
	expected := entity.User{
		ID:       uuid.NewV4(),
		Username: gophtest.Username,
	}

	rows := pgxmock.NewRows([]string{"user_id", "username"}).
		AddRow(expected.ID.String(), expected.Username)

	m := newPoolMock(t)
	m.ExpectQuery("SELECT user_id, username FROM users").
		WithArgs(gophtest.Username, gophtest.SecurityKey).
		WillReturnRows(rows)

	sat := newUsersRepo(t, m)
	rv, err := sat.Verify(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.NoError(t, err)
	require.Equal(t, expected, rv)
	require.NoError(t, m.ExpectationsWereMet())
}

func TestVerifyFailsOnBadCredentials(t *testing.T) {
	rows := pgxmock.NewRows([]string{"user_id", "username"})

	m := newPoolMock(t)
	m.ExpectQuery("SELECT user_id, username FROM users").
		WithArgs(gophtest.Username, gophtest.SecurityKey).
		WillReturnRows(rows)

	sat := newUsersRepo(t, m)
	_, err := sat.Verify(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.ErrorIs(t, err, entity.ErrInvalidCredentials)
	require.NoError(t, m.ExpectationsWereMet())
}

func TestVerifyFailsOnUnexpectedError(t *testing.T) {
	m := newPoolMock(t)
	m.ExpectQuery("SELECT user_id, username FROM users").
		WithArgs(gophtest.Username, gophtest.SecurityKey).
		WillReturnError(gophtest.ErrUnexpected)

	sat := newUsersRepo(t, m)
	_, err := sat.Verify(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.ErrorIs(t, err, gophtest.ErrUnexpected)
	require.NoError(t, m.ExpectationsWereMet())
}
