package repo_test

import (
	"context"
	"testing"

	"github.com/alkurbatov/goph-keeper/internal/keeper/entity"
	"github.com/alkurbatov/goph-keeper/internal/keeper/infra/postgres"
	"github.com/alkurbatov/goph-keeper/internal/libraries/gophtest"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
	"github.com/pashagolub/pgxmock/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func doGetSecret(
	t *testing.T,
	owner, id uuid.UUID,
	m pgxmock.PgxPoolIface,
) (*entity.Secret, error) {
	t.Helper()

	sat := newTestRepos(t, m).Secrets
	secret, err := sat.Get(context.Background(), owner, id)

	require.NoError(t, m.ExpectationsWereMet())

	return secret, err
}

func doDeleteSecret(t *testing.T, owner, id uuid.UUID, m pgxmock.PgxPoolIface) error {
	t.Helper()

	sat := newTestRepos(t, m).Secrets
	err := sat.Delete(context.Background(), owner, id)

	require.NoError(t, m.ExpectationsWereMet())

	return err
}

func TestCreateSecret(t *testing.T) {
	owner := uuid.NewV4()
	expected := uuid.NewV4()

	rows := pgxmock.NewRows([]string{"id"}).
		AddRow(expected.String())

	m := newPoolMock(t)
	m.ExpectBeginTx(postgres.DefaultTxOptions)
	m.ExpectQuery("INSERT INTO secrets").
		WithArgs(
			owner,
			gophtest.SecretName,
			goph.DataKind_TEXT,
			[]byte(gophtest.Metadata),
			[]byte(gophtest.TextData),
		).
		WillReturnRows(rows)
	m.ExpectCommit()

	sat := newTestRepos(t, m).Secrets
	id, err := sat.Create(
		context.Background(),
		owner,
		gophtest.SecretName,
		goph.DataKind_TEXT,
		[]byte(gophtest.Metadata),
		[]byte(gophtest.TextData),
	)

	require.NoError(t, err)
	require.Equal(t, expected, id)
	require.NoError(t, m.ExpectationsWereMet())
}

func TestCreateSecretFailure(t *testing.T) {
	tt := []struct {
		name     string
		err      error
		expected error
	}{
		{
			name:     "Create secret fails if secret exists",
			err:      errUniqueViolation,
			expected: entity.ErrSecretExists,
		},
		{
			name:     "Create secret fails on unexpected error",
			err:      gophtest.ErrUnexpected,
			expected: gophtest.ErrUnexpected,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			owner := uuid.NewV4()

			m := newPoolMock(t)
			m.ExpectBeginTx(postgres.DefaultTxOptions)
			m.ExpectQuery("INSERT").
				WithArgs(
					owner,
					gophtest.SecretName,
					goph.DataKind_TEXT,
					[]byte(gophtest.Metadata),
					[]byte(gophtest.TextData),
				).
				WillReturnError(tc.err)
			m.ExpectRollback()

			sat := newTestRepos(t, m).Secrets
			_, err := sat.Create(
				context.Background(),
				owner,
				gophtest.SecretName,
				goph.DataKind_TEXT,
				[]byte(gophtest.Metadata),
				[]byte(gophtest.TextData),
			)

			require.ErrorIs(t, err, tc.expected)
			require.NoError(t, m.ExpectationsWereMet())
		})
	}
}

func TestListSecrets(t *testing.T) {
	tt := []struct {
		name string
		rows [][]any
	}{
		{
			name: "List secrets of a user",
			rows: [][]any{
				{uuid.NewV4().String(), gophtest.SecretName, goph.DataKind_TEXT, []byte("xxx")},
				{uuid.NewV4().String(), gophtest.SecretName + "ex", goph.DataKind_BINARY, []byte{}},
			},
		},
		{
			name: "List secrets returns empty list",
			rows: [][]any{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			owner := uuid.NewV4()
			rows := pgxmock.NewRows([]string{"secret_id", "name", "kind", "metadata"})

			for _, row := range tc.rows {
				rows.AddRow(row...)
			}

			m := newPoolMock(t)
			m.ExpectQuery("SELECT secret_id, name, kind, metadata FROM secrets").
				WithArgs(owner).
				WillReturnRows(rows)

			sat := newTestRepos(t, m).Secrets
			secrets, err := sat.List(context.Background(), owner)

			require.NoError(t, err)
			require.Len(t, secrets, len(tc.rows))
			require.NoError(t, m.ExpectationsWereMet())
		})
	}
}

func TestListSecretsOnFailure(t *testing.T) {
	owner := uuid.NewV4()

	m := newPoolMock(t)
	m.ExpectQuery("SELECT").
		WithArgs(owner).
		WillReturnError(gophtest.ErrUnexpected)

	sat := newTestRepos(t, m).Secrets
	_, err := sat.List(context.Background(), owner)

	require.Error(t, err)
	require.NoError(t, m.ExpectationsWereMet())
}

func TestGetSecret(t *testing.T) {
	owner := uuid.NewV4()

	expected := &entity.Secret{
		ID:       uuid.NewV4(),
		Name:     gophtest.SecretName,
		Kind:     goph.DataKind_TEXT,
		Metadata: []byte(gophtest.Metadata),
		Data:     []byte(gophtest.TextData),
	}

	rows := pgxmock.NewRows([]string{"secret_id", "name", "kind", "metadata", "data"}).
		AddRow(expected.ID.String(), expected.Name, expected.Kind, expected.Metadata, expected.Data)

	m := newPoolMock(t)
	m.ExpectQuery("SELECT secret_id, name, kind, metadata, data FROM secrets").
		WithArgs(expected.ID, owner).
		WillReturnRows(rows)

	secret, err := doGetSecret(t, owner, expected.ID, m)

	require.NoError(t, err)
	require.Equal(t, expected, secret)
}

func TestGetUnexistingSecret(t *testing.T) {
	rows := pgxmock.NewRows([]string{"secret_id", "name", "kind", "metadata", "data"})

	owner := uuid.NewV4()
	id := uuid.NewV4()

	m := newPoolMock(t)
	m.ExpectQuery("SELECT").
		WithArgs(id, owner).
		WillReturnRows(rows)

	_, err := doGetSecret(t, owner, id, m)

	require.ErrorIs(t, err, entity.ErrSecretNotFound)
}

func TestGetSecretOnFailure(t *testing.T) {
	owner := uuid.NewV4()
	id := uuid.NewV4()

	m := newPoolMock(t)
	m.ExpectQuery("SELECT").
		WithArgs(id, owner).
		WillReturnError(gophtest.ErrUnexpected)

	_, err := doGetSecret(t, owner, id, m)

	require.ErrorIs(t, err, gophtest.ErrUnexpected)
}

func TestDeleteSecret(t *testing.T) {
	owner := uuid.NewV4()
	id := uuid.NewV4()

	m := newPoolMock(t)
	m.ExpectBeginTx(postgres.DefaultTxOptions)
	m.ExpectExec("DELETE FROM secrets").
		WithArgs(id, owner).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	m.ExpectCommit()

	err := doDeleteSecret(t, owner, id, m)

	require.NoError(t, err)
}

func TestDeleteUnexistingSecret(t *testing.T) {
	owner := uuid.NewV4()
	id := uuid.NewV4()

	m := newPoolMock(t)
	m.ExpectBeginTx(postgres.DefaultTxOptions)
	m.ExpectExec("DELETE").
		WithArgs(id, owner).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))
	m.ExpectRollback()

	err := doDeleteSecret(t, owner, id, m)

	require.ErrorIs(t, err, entity.ErrSecretNotFound)
}

func TestDeleteSecretOnFailure(t *testing.T) {
	owner := uuid.NewV4()
	id := uuid.NewV4()

	m := newPoolMock(t)
	m.ExpectBeginTx(postgres.DefaultTxOptions)
	m.ExpectExec("DELETE").
		WithArgs(id, owner).
		WillReturnError(gophtest.ErrUnexpected)
	m.ExpectRollback()

	err := doDeleteSecret(t, owner, id, m)

	require.ErrorIs(t, err, gophtest.ErrUnexpected)
}
