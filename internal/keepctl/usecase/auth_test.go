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
	key := newTestKey()

	m := &repo.AuthRepoMock{}
	m.On(
		"Login",
		mock.Anything,
		gophtest.Username,
		key.Hash(),
	).
		Return(gophtest.AccessToken, nil)

	sat := usecase.NewAuthUseCase(m)
	token, err := sat.Login(context.Background(), gophtest.Username, key)

	require.NoError(t, err)
	require.Equal(t, gophtest.AccessToken, token)
	m.AssertExpectations(t)
}

func TestLoginOnRepoFailure(t *testing.T) {
	key := newTestKey()

	m := &repo.AuthRepoMock{}
	m.On(
		"Login",
		mock.Anything,
		gophtest.Username,
		key.Hash(),
	).
		Return("", gophtest.ErrUnexpected)

	sat := usecase.NewAuthUseCase(m)
	_, err := sat.Login(context.Background(), gophtest.Username, key)

	require.Error(t, err)
	m.AssertExpectations(t)
}
