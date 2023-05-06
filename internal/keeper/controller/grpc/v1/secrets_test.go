package v1_test

import (
	"context"
	"strings"
	"testing"

	v1 "github.com/alkurbatov/goph-keeper/internal/keeper/controller/grpc/v1"
	"github.com/alkurbatov/goph-keeper/internal/keeper/entity"
	"github.com/alkurbatov/goph-keeper/internal/keeper/usecase"
	"github.com/alkurbatov/goph-keeper/internal/libraries/gophtest"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestCreateSecret(t *testing.T) {
	tt := []struct {
		name       string
		secretName string
		metadata   []byte
		data       []byte
	}{
		{
			name:       "Create secret",
			secretName: gophtest.SecretName,
			metadata:   []byte(gophtest.Metadata),
			data:       []byte(gophtest.TextData),
		},
		{
			name:       "Create secret without metadata",
			secretName: gophtest.Username,
			metadata:   nil,
			data:       []byte(gophtest.TextData),
		},
		{
			name:       "Create secret of maximum size",
			secretName: strings.Repeat("#", v1.DefaultMaxSecretNameLength),
			metadata:   []byte(strings.Repeat("#", v1.DefaultMetadataLimit)),
			data:       []byte(strings.Repeat("#", v1.DefaultDataLimit)),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			expected := uuid.NewV4()

			m := newUseCasesMock()
			m.Secrets.(*usecase.SecretsUseCaseMock).On(
				"Create",
				mock.Anything,
				mock.AnythingOfType("uuid.UUID"),
				tc.secretName,
				goph.DataKind_BINARY,
				tc.metadata,
				tc.data,
			).
				Return(expected, nil)

			conn := createTestServerWithFakeAuth(t, m)

			req := &goph.CreateSecretRequest{
				Name:     tc.secretName,
				Kind:     goph.DataKind_BINARY,
				Metadata: tc.metadata,
				Data:     tc.data,
			}

			client := goph.NewSecretsClient(conn)
			resp, err := client.Create(context.Background(), req)

			require.NoError(t, err)
			require.Equal(t, expected.String(), resp.GetId())
			m.Secrets.(*usecase.SecretsUseCaseMock).AssertExpectations(t)
		})
	}
}

func TestCreateSecretWithBadRequest(t *testing.T) {
	tt := []struct {
		name       string
		secretName string
		metadata   []byte
		data       []byte
	}{
		{
			name:       "Create secret fails if secret name is empty",
			secretName: "",
			metadata:   []byte(gophtest.Metadata),
			data:       []byte(gophtest.TextData),
		},
		{
			name:       "Create secret fails if secret name is too long",
			secretName: strings.Repeat("#", v1.DefaultMaxSecretNameLength+1),
			metadata:   []byte(gophtest.Metadata),
			data:       []byte(gophtest.TextData),
		},
		{
			name:       "Create secret fails if metadata is too long",
			secretName: gophtest.Username,
			metadata:   []byte(strings.Repeat("#", v1.DefaultMetadataLimit+1)),
			data:       make([]byte, 0),
		},
		{
			name:       "Create secret fails if data is empty",
			secretName: gophtest.Username,
			metadata:   []byte(gophtest.Metadata),
			data:       make([]byte, 0),
		},
		{
			name:       "Create secret fails if data is too long",
			secretName: gophtest.Username,
			metadata:   []byte(gophtest.Metadata),
			data:       []byte(strings.Repeat("#", v1.DefaultDataLimit+1)),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			conn := createTestServer(t, newUseCasesMock())

			req := &goph.CreateSecretRequest{
				Name:     tc.secretName,
				Kind:     goph.DataKind_BINARY,
				Metadata: tc.metadata,
				Data:     tc.data,
			}

			client := goph.NewSecretsClient(conn)
			_, err := client.Create(context.Background(), req)

			requireEqualCode(t, codes.InvalidArgument, err)
		})
	}
}

func TestCreateSecretFailsIfNoUserInfo(t *testing.T) {
	conn := createTestServer(t, newUseCasesMock())

	req := &goph.CreateSecretRequest{
		Name:     gophtest.SecretName,
		Kind:     goph.DataKind_BINARY,
		Metadata: []byte(gophtest.Metadata),
		Data:     []byte(gophtest.TextData),
	}

	client := goph.NewSecretsClient(conn)
	_, err := client.Create(context.Background(), req)

	requireEqualCode(t, codes.Unauthenticated, err)
}

func TestCreateServerOnUseCaseFailure(t *testing.T) {
	tt := []struct {
		name     string
		err      error
		expected codes.Code
	}{
		{
			name:     "Create secret fails if secret already exists",
			err:      entity.ErrSecretExists,
			expected: codes.AlreadyExists,
		},
		{
			name:     "Create secret fails if use case fails unexpectedly",
			err:      gophtest.ErrUnexpected,
			expected: codes.Internal,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := newUseCasesMock()
			m.Secrets.(*usecase.SecretsUseCaseMock).On(
				"Create",
				mock.Anything,
				mock.AnythingOfType("uuid.UUID"),
				gophtest.SecretName,
				goph.DataKind_BINARY,
				[]byte(gophtest.Metadata),
				[]byte(gophtest.TextData),
			).
				Return(uuid.UUID{}, tc.err)

			conn := createTestServerWithFakeAuth(t, m)

			req := &goph.CreateSecretRequest{
				Name:     gophtest.SecretName,
				Kind:     goph.DataKind_BINARY,
				Metadata: []byte(gophtest.Metadata),
				Data:     []byte(gophtest.TextData),
			}

			client := goph.NewSecretsClient(conn)
			_, err := client.Create(context.Background(), req)

			requireEqualCode(t, tc.expected, err)
			m.Secrets.(*usecase.SecretsUseCaseMock).AssertExpectations(t)
		})
	}
}
