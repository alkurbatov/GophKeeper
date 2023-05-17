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

func TestRegister(t *testing.T) {
	key := newTestKey()

	m := &repo.UsersRepoMock{}
	m.On(
		"Register",
		mock.Anything,
		gophtest.Username,
		key.Hash(),
	).
		Return(gophtest.AccessToken, nil)

	sat := usecase.NewUsersUseCase(m)
	token, err := sat.Register(context.Background(), gophtest.Username, key)

	require.NoError(t, err)
	require.Equal(t, gophtest.AccessToken, token)
	m.AssertExpectations(t)
}

func TestRegisterOnRepoFailure(t *testing.T) {
	key := newTestKey()

	m := &repo.UsersRepoMock{}
	m.On(
		"Register",
		mock.Anything,
		gophtest.Username,
		key.Hash(),
	).
		Return("", gophtest.ErrUnexpected)

	sat := usecase.NewUsersUseCase(m)
	_, err := sat.Register(context.Background(), gophtest.Username, key)

	require.Error(t, err)
	m.AssertExpectations(t)
}
