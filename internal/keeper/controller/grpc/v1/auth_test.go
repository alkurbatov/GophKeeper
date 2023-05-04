package v1_test

import (
	"context"
	"strings"
	"testing"

	"github.com/alkurbatov/goph-keeper/internal/keeper/entity"
	"github.com/alkurbatov/goph-keeper/internal/keeper/usecase"
	"github.com/alkurbatov/goph-keeper/internal/libraries/gophtest"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestLoginUser(t *testing.T) {
	m := newUseCasesMock()
	m.Auth.(*usecase.AuthUseCaseMock).On(
		"Login",
		mock.Anything,
		gophtest.Username,
		gophtest.SecurityKey,
	).
		Return(entity.AccessToken(gophtest.AccessToken), nil)

	conn := createTestServer(t, m)

	req := &goph.LoginRequest{
		Username:    gophtest.Username,
		SecurityKey: gophtest.SecurityKey,
	}

	client := goph.NewAuthClient(conn)
	resp, err := client.Login(context.Background(), req)

	require.NoError(t, err)
	require.Equal(t, gophtest.AccessToken, resp.AccessToken)
	m.Auth.(*usecase.AuthUseCaseMock).AssertExpectations(t)
}

func TestLoginWithBadRequest(t *testing.T) {
	tt := []struct {
		name     string
		username string
		key      string
	}{
		{
			name:     "Login fails if username is empty",
			username: "",
			key:      gophtest.SecurityKey,
		},
		{
			name:     "Login fails if security key is empty",
			username: gophtest.Username,
			key:      "",
		},
		{
			name:     "Login fails if username is too long",
			username: strings.Repeat("#", 129),
			key:      gophtest.SecurityKey,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			conn := createTestServer(t, newUseCasesMock())

			req := &goph.LoginRequest{
				Username:    tc.username,
				SecurityKey: tc.key,
			}

			client := goph.NewAuthClient(conn)
			_, err := client.Login(context.Background(), req)

			requireEqualCode(t, codes.InvalidArgument, err)
		})
	}
}

func TestLoginOnUseCaseFailure(t *testing.T) {
	tt := []struct {
		name       string
		useCaseErr error
		expected   codes.Code
	}{
		{
			name:       "Login fails on invalid credentials",
			useCaseErr: entity.ErrInvalidCredentials,
			expected:   codes.Unauthenticated,
		},
		{
			name:       "Login fails if something bad happened",
			useCaseErr: gophtest.ErrUnexpected,
			expected:   codes.Internal,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := newUseCasesMock()
			m.Auth.(*usecase.AuthUseCaseMock).On(
				"Login",
				mock.Anything,
				gophtest.Username,
				gophtest.SecurityKey,
			).
				Return(entity.AccessToken(""), tc.useCaseErr)

			conn := createTestServer(t, m)

			req := &goph.LoginRequest{
				Username:    gophtest.Username,
				SecurityKey: gophtest.SecurityKey,
			}

			client := goph.NewAuthClient(conn)
			_, err := client.Login(context.Background(), req)

			requireEqualCode(t, tc.expected, err)
			m.Auth.(*usecase.AuthUseCaseMock).AssertExpectations(t)
		})
	}
}
