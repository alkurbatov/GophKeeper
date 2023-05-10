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
	"github.com/gkampitakis/go-snaps/snaps"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func doListSecrets(
	t *testing.T,
	mockRV []entity.Secret,
	mockErr error,
) (*goph.ListSecretsResponse, error) {
	t.Helper()

	m := newUseCasesMock()
	m.Secrets.(*usecase.SecretsUseCaseMock).On(
		"List",
		mock.Anything,
		mock.AnythingOfType("uuid.UUID"),
	).
		Return(mockRV, mockErr)

	conn := createTestServerWithFakeAuth(t, m)
	req := &goph.ListSecretsRequest{}

	client := goph.NewSecretsClient(conn)
	rv, err := client.List(context.Background(), req)

	m.Secrets.(*usecase.SecretsUseCaseMock).AssertExpectations(t)

	return rv, err
}

func doGetSecret(
	t *testing.T,
	mockRV *entity.Secret,
	mockErr error,
) (*goph.GetSecretResponse, error) {
	t.Helper()

	m := newUseCasesMock()
	m.Secrets.(*usecase.SecretsUseCaseMock).On(
		"Get",
		mock.Anything,
		mock.AnythingOfType("uuid.UUID"),
		mock.AnythingOfType("uuid.UUID"),
	).
		Return(mockRV, mockErr)

	conn := createTestServerWithFakeAuth(t, m)
	req := &goph.GetSecretRequest{Id: uuid.NewV4().String()}

	client := goph.NewSecretsClient(conn)
	rv, err := client.Get(context.Background(), req)

	m.Secrets.(*usecase.SecretsUseCaseMock).AssertExpectations(t)

	return rv, err
}

func doDeleteSecret(
	t *testing.T,
	mockErr error,
) (*goph.DeleteSecretResponse, error) {
	t.Helper()

	m := newUseCasesMock()
	m.Secrets.(*usecase.SecretsUseCaseMock).On(
		"Delete",
		mock.Anything,
		mock.AnythingOfType("uuid.UUID"),
		mock.AnythingOfType("uuid.UUID"),
	).
		Return(mockErr)

	conn := createTestServerWithFakeAuth(t, m)
	req := &goph.DeleteSecretRequest{Id: uuid.NewV4().String()}

	client := goph.NewSecretsClient(conn)
	rv, err := client.Delete(context.Background(), req)

	m.Secrets.(*usecase.SecretsUseCaseMock).AssertExpectations(t)

	return rv, err
}

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

func TestListSecrets(t *testing.T) {
	tt := []struct {
		name    string
		secrets []entity.Secret
	}{
		{
			name: "List secrets of a user",
			secrets: []entity.Secret{
				{
					ID:   gophtest.CreateUUID(t, "7728154c-9400-4f1b-a2a3-01deb83ece05"),
					Name: gophtest.SecretName,
					Kind: goph.DataKind_BINARY,
				},
				{
					ID:       gophtest.CreateUUID(t, "df566e25-43a5-4c34-9123-3931fb809b45"),
					Name:     gophtest.SecretName + "ex",
					Kind:     goph.DataKind_TEXT,
					Metadata: []byte(gophtest.Metadata),
				},
			},
		},
		{
			name:    "List secrets when user has no secrets",
			secrets: []entity.Secret{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rv, err := doListSecrets(t, tc.secrets, nil)

			require.NoError(t, err)
			snaps.MatchSnapshot(t, rv.GetSecrets())
		})
	}
}

func TestListSecretsFailsIfNoUserInfo(t *testing.T) {
	conn := createTestServer(t, newUseCasesMock())

	req := &goph.ListSecretsRequest{}

	client := goph.NewSecretsClient(conn)
	_, err := client.List(context.Background(), req)

	requireEqualCode(t, codes.Unauthenticated, err)
}

func TestListSecretsFailsOnUseCaseFailure(t *testing.T) {
	_, err := doListSecrets(t, nil, gophtest.ErrUnexpected)

	requireEqualCode(t, codes.Internal, err)
}

func TestGetSecret(t *testing.T) {
	tt := []struct {
		name   string
		secret *entity.Secret
	}{
		{
			name: "Get secret",
			secret: &entity.Secret{
				ID:       gophtest.CreateUUID(t, "df566e25-43a5-4c34-9123-3931fb809b45"),
				Name:     gophtest.SecretName,
				Kind:     goph.DataKind_TEXT,
				Metadata: []byte(gophtest.Metadata),
				Data:     []byte(gophtest.TextData),
			},
		},
		{
			name: "Get secret without metadata",
			secret: &entity.Secret{
				ID:       gophtest.CreateUUID(t, "df566e25-43a5-4c34-9123-3931fb809b45"),
				Name:     gophtest.SecretName,
				Kind:     goph.DataKind_TEXT,
				Metadata: []byte(gophtest.Metadata),
				Data:     []byte(gophtest.TextData),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := doGetSecret(t, tc.secret, nil)

			require.NoError(t, err)
			snaps.MatchSnapshot(t, resp.GetSecret())
			require.Equal(t, tc.secret.Data, resp.GetData())
		})
	}
}

func TestGetSecretOnBadRequest(t *testing.T) {
	conn := createTestServerWithFakeAuth(t, newUseCasesMock())

	req := &goph.GetSecretRequest{Id: "xxx"}

	client := goph.NewSecretsClient(conn)
	_, err := client.Get(context.Background(), req)

	requireEqualCode(t, codes.InvalidArgument, err)
}

func TestGetSecretFailsIfNoUserInfo(t *testing.T) {
	conn := createTestServer(t, newUseCasesMock())

	req := &goph.GetSecretRequest{Id: uuid.NewV4().String()}

	client := goph.NewSecretsClient(conn)
	_, err := client.Get(context.Background(), req)

	requireEqualCode(t, codes.Unauthenticated, err)
}

func TestGetSecretFailsOnUseCaseFailure(t *testing.T) {
	_, err := doGetSecret(t, nil, gophtest.ErrUnexpected)

	requireEqualCode(t, codes.Internal, err)
}

func TestDeleteSecret(t *testing.T) {
	_, err := doDeleteSecret(t, nil)

	require.NoError(t, err)
}

func TestDeleteSecretOnBadRequest(t *testing.T) {
	conn := createTestServerWithFakeAuth(t, newUseCasesMock())

	req := &goph.DeleteSecretRequest{Id: "xxx"}

	client := goph.NewSecretsClient(conn)
	_, err := client.Delete(context.Background(), req)

	requireEqualCode(t, codes.InvalidArgument, err)
}

func TestDeleteSecretFailsIfNoUserInfo(t *testing.T) {
	conn := createTestServer(t, newUseCasesMock())

	req := &goph.DeleteSecretRequest{Id: uuid.NewV4().String()}

	client := goph.NewSecretsClient(conn)
	_, err := client.Delete(context.Background(), req)

	requireEqualCode(t, codes.Unauthenticated, err)
}

func TestDeleteSecretFailsOnUseCaseFailure(t *testing.T) {
	_, err := doDeleteSecret(t, gophtest.ErrUnexpected)

	requireEqualCode(t, codes.Internal, err)
}
