package repo_test

import (
	"context"
	"testing"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/repo"
	"github.com/alkurbatov/goph-keeper/internal/libraries/gophtest"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func newLoginRequest() *goph.LoginRequest {
	return &goph.LoginRequest{
		Username:    gophtest.Username,
		SecurityKey: gophtest.SecurityKey,
	}
}

func TestLogin(t *testing.T) {
	resp := &goph.LoginResponse{
		AccessToken: gophtest.AccessToken,
	}

	m := &goph.AuthClientMock{}
	m.On(
		"Login",
		mock.Anything,
		newLoginRequest(),
		mock.Anything,
	).
		Return(resp, nil)

	sat := repo.NewAuthRepo(m)
	token, err := sat.Login(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.NoError(t, err)
	require.Equal(t, gophtest.AccessToken, token)
	m.AssertExpectations(t)
}

func TestLoginOnOperationFailure(t *testing.T) {
	m := &goph.AuthClientMock{}
	m.On(
		"Login",
		mock.Anything,
		newLoginRequest(),
		mock.Anything,
	).
		Return(nil, grpc.ErrServerStopped)

	sat := repo.NewAuthRepo(m)
	_, err := sat.Login(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.Error(t, err)
	m.AssertExpectations(t)
}
