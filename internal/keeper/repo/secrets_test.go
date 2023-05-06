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
			goph.DataKind_TEXT.String(),
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
			m.ExpectQuery("INSERT INTO secrets").
				WithArgs(
					owner,
					gophtest.SecretName,
					goph.DataKind_TEXT.String(),
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
