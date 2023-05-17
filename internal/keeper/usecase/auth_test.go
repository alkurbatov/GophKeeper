package usecase_test

import (
	"context"
	"testing"

	"github.com/alkurbatov/goph-keeper/internal/keeper/entity"
	"github.com/alkurbatov/goph-keeper/internal/keeper/repo"
	"github.com/alkurbatov/goph-keeper/internal/keeper/usecase"
	"github.com/alkurbatov/goph-keeper/internal/libraries/gophtest"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func doLogin(t *testing.T, repoErr error) (entity.AccessToken, error) {
	t.Helper()

	m := &repo.UsersRepoMock{}
	m.On(
		"Verify",
		mock.Anything,
		gophtest.Username,
		gophtest.SecurityKey,
	).
		Return(entity.User{ID: uuid.NewV4(), Username: gophtest.Username}, repoErr)

	sat := usecase.NewAuthUseCase(gophtest.Secret, m)
	accessToken, err := sat.Login(context.Background(), gophtest.Username, gophtest.SecurityKey)

	m.AssertExpectations(t)

	return accessToken, err
}

func TestLogin(t *testing.T) {
	token, err := doLogin(t, nil)

	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestLoginOnBadCredentials(t *testing.T) {
	_, err := doLogin(t, entity.ErrInvalidCredentials)

	require.ErrorIs(t, err, entity.ErrInvalidCredentials)
}
