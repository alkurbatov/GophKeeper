package repo_test

import (
	"context"
	"testing"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/repo"
	"github.com/alkurbatov/goph-keeper/internal/libraries/gophtest"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
	"github.com/gkampitakis/go-snaps/snaps"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func doCreateSecret(
	t *testing.T,
	mockRV *goph.CreateSecretResponse,
	mockErr error,
) (uuid.UUID, error) {
	t.Helper()

	req := &goph.CreateSecretRequest{
		Name:     gophtest.SecretName,
		Metadata: []byte(gophtest.Metadata),
		Kind:     goph.DataKind_TEXT,
		Data:     []byte(gophtest.TextData),
	}

	m := &goph.SecretsClientMock{}
	m.On(
		"Create",
		mock.Anything,
		req,
		mock.Anything,
	).
		Return(mockRV, mockErr)

	sat := repo.NewSecretsRepo(m)
	rv, err := sat.Push(
		context.Background(),
		gophtest.AccessToken,
		gophtest.SecretName,
		goph.DataKind_TEXT,
		[]byte(gophtest.Metadata),
		[]byte(gophtest.TextData),
	)

	m.AssertExpectations(t)

	return rv, err
}

func doListSecrets(
	t *testing.T,
	mockRV *goph.ListSecretsResponse,
	mockErr error,
) ([]*goph.Secret, error) {
	t.Helper()

	req := &goph.ListSecretsRequest{}

	m := &goph.SecretsClientMock{}
	m.On(
		"List",
		mock.Anything,
		req,
		mock.Anything,
	).
		Return(mockRV, mockErr)

	sat := repo.NewSecretsRepo(m)
	rv, err := sat.List(context.Background(), gophtest.AccessToken)

	m.AssertExpectations(t)

	return rv, err
}

func TestCreateSecret(t *testing.T) {
	expected := uuid.NewV4()
	resp := &goph.CreateSecretResponse{
		Id: expected.String(),
	}

	id, err := doCreateSecret(t, resp, nil)

	require.NoError(t, err)
	require.Equal(t, expected, id)
}

func TestCreateSecretOnOperationFailure(t *testing.T) {
	_, err := doCreateSecret(t, nil, gophtest.ErrUnexpected)

	require.Error(t, err)
}

func TestListSecrets(t *testing.T) {
	tt := []struct {
		name    string
		secrets []*goph.Secret
	}{
		{
			name: "List secrets of a user",
			secrets: []*goph.Secret{
				{
					Id:       gophtest.CreateUUID(t, "df566e25-43a5-4c34-9123-3931fb809b45").String(),
					Name:     gophtest.SecretName,
					Kind:     goph.DataKind_TEXT,
					Metadata: []byte(gophtest.Metadata),
				},
			},
		},
		{
			name:    "List secrets of a user who has no secrets",
			secrets: []*goph.Secret{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resp := &goph.ListSecretsResponse{
				Secrets: tc.secrets,
			}

			rv, err := doListSecrets(t, resp, nil)

			require.NoError(t, err)
			snaps.MatchSnapshot(t, rv)
		})
	}
}

func TestListSecretsOnOperationFailure(t *testing.T) {
	_, err := doListSecrets(t, nil, gophtest.ErrUnexpected)

	require.Error(t, err)
}
