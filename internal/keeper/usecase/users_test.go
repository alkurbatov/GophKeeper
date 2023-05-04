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

func TestRegister(t *testing.T) {
	m := &repo.UsersRepoMock{}
	m.On(
		"Register",
		mock.Anything,
		gophtest.Username,
		gophtest.SecurityKey,
	).
		Return(uuid.NewV4(), nil)

	sat := usecase.NewUsersUseCase(gophtest.Secret, m)
	token, err := sat.Register(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.NoError(t, err)
	require.NotEmpty(t, token)
	m.AssertExpectations(t)
}

func TestRegisterFailsIfUserExists(t *testing.T) {
	m := &repo.UsersRepoMock{}
	m.On(
		"Register",
		mock.Anything,
		gophtest.Username,
		gophtest.SecurityKey,
	).
		Return(uuid.UUID{}, entity.ErrUserExists)

	sat := usecase.NewUsersUseCase(gophtest.Secret, m)
	_, err := sat.Register(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.Error(t, err)
	m.AssertExpectations(t)
}
