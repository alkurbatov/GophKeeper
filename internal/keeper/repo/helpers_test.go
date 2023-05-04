package repo_test

import (
	"testing"

	"github.com/alkurbatov/goph-keeper/internal/keeper/infra/logger"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/require"
)

func createTestLogger(t *testing.T) *logger.Logger {
	t.Helper()

	log, err := logger.New("info")
	require.NoError(t, err)

	return log
}

func newPoolMock(t *testing.T) pgxmock.PgxPoolIface {
	t.Helper()

	m, err := pgxmock.NewPool()
	require.NoError(t, err)

	t.Cleanup(m.Close)

	return m
}
