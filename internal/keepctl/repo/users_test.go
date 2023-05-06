package repo_test

import (
	"context"
	"testing"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/repo"
	"github.com/alkurbatov/goph-keeper/internal/libraries/gophtest"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newRegisterUserRequest() *goph.RegisterUserRequest {
	return &goph.RegisterUserRequest{
		Username:    gophtest.Username,
		SecurityKey: gophtest.SecurityKey,
	}
}

func TestRegister(t *testing.T) {
	resp := &goph.RegisterUserResponse{
		AccessToken: gophtest.AccessToken,
	}

	m := &goph.UsersClientMock{}
	m.On(
		"Register",
		mock.Anything,
		newRegisterUserRequest(),
		mock.Anything,
	).
		Return(resp, nil)

	sat := repo.NewUsersRepo(m)
	token, err := sat.Register(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.NoError(t, err)
	require.Equal(t, gophtest.AccessToken, token)
	m.AssertExpectations(t)
}

func TestRegisterOnOperationFailure(t *testing.T) {
	m := &goph.UsersClientMock{}
	m.On(
		"Register",
		mock.Anything,
		newRegisterUserRequest(),
		mock.Anything,
	).
		Return(nil, gophtest.ErrUnexpected)

	sat := repo.NewUsersRepo(m)
	_, err := sat.Register(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.Error(t, err)
	m.AssertExpectations(t)
}
