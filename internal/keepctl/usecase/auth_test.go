package usecase_test

import (
	"context"
	"testing"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/repo"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/usecase"
	"github.com/alkurbatov/goph-keeper/internal/libraries/gophtest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	m := &repo.AuthRepoMock{}
	m.On(
		"Login",
		mock.Anything,
		gophtest.Username,
		mock.AnythingOfType("string"),
	).
		Return(gophtest.AccessToken, nil)

	sat := usecase.NewAuthUseCase(m)
	token, err := sat.Login(context.Background(), gophtest.Username, gophtest.Password)

	require.NoError(t, err)
	require.Equal(t, gophtest.AccessToken, token)
	m.AssertExpectations(t)
}

func TestLoginOnRepoFailure(t *testing.T) {
	m := &repo.AuthRepoMock{}
	m.On(
		"Login",
		mock.Anything,
		gophtest.Username,
		mock.AnythingOfType("string"),
	).
		Return("", gophtest.ErrUnexpected)

	sat := usecase.NewAuthUseCase(m)
	_, err := sat.Login(context.Background(), gophtest.Username, gophtest.Password)

	require.Error(t, err)
	m.AssertExpectations(t)
}
