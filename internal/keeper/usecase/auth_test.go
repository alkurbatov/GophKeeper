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

func TestLogin(t *testing.T) {
	m := &repo.UsersRepoMock{}
	m.On(
		"Verify",
		mock.Anything,
		gophtest.Username,
		gophtest.SecurityKey,
	).
		Return(entity.User{ID: uuid.NewV4(), Username: gophtest.Username}, nil)

	sat := usecase.NewAuthUseCase(gophtest.Secret, m)
	accessToken, err := sat.Login(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.NoError(t, err)
	require.NotEmpty(t, accessToken)
	m.AssertExpectations(t)
}

func TestLoginOnBadCredentials(t *testing.T) {
	m := &repo.UsersRepoMock{}
	m.On(
		"Verify",
		mock.Anything,
		gophtest.Username,
		gophtest.SecurityKey,
	).
		Return(entity.User{}, entity.ErrInvalidCredentials)

	sat := usecase.NewAuthUseCase(gophtest.Secret, m)
	_, err := sat.Login(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.ErrorIs(t, err, entity.ErrInvalidCredentials)
	m.AssertExpectations(t)
}
