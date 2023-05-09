package usecase_test

import (
	"context"
	"testing"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/repo"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/usecase"
	"github.com/alkurbatov/goph-keeper/internal/libraries/gophtest"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
	"github.com/gkampitakis/go-snaps/snaps"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func doPushText(t *testing.T, mockRV uuid.UUID, mockErr error) (uuid.UUID, error) {
	t.Helper()

	m := &repo.SecretsRepoMock{}
	m.On(
		"Push",
		mock.Anything,
		gophtest.AccessToken,
		gophtest.SecretName,
		goph.DataKind_TEXT,
		mock.AnythingOfType("[]uint8"),
		mock.AnythingOfType("[]uint8"),
	).
		Return(mockRV, mockErr)

	sat := usecase.NewSecretsUseCase(newTestKey(), m)
	id, err := sat.PushText(
		context.Background(),
		gophtest.AccessToken,
		gophtest.SecretName,
		gophtest.Metadata,
		gophtest.TextData,
	)

	m.AssertExpectations(t)

	return id, err
}

func doList(t *testing.T, mockRV []*goph.Secret, mockErr error) ([]*goph.Secret, error) {
	t.Helper()

	m := &repo.SecretsRepoMock{}
	m.On(
		"List",
		mock.Anything,
		gophtest.AccessToken,
	).
		Return(mockRV, mockErr)

	sat := usecase.NewSecretsUseCase(newTestKey(), m)
	data, err := sat.List(
		context.Background(),
		gophtest.AccessToken,
	)

	m.AssertExpectations(t)

	return data, err
}

func doDelete(t *testing.T, mockErr error) error {
	t.Helper()

	id := uuid.NewV4()

	m := &repo.SecretsRepoMock{}
	m.On(
		"Delete",
		mock.Anything,
		gophtest.AccessToken,
		id,
	).
		Return(mockErr)

	sat := usecase.NewSecretsUseCase(newTestKey(), m)
	err := sat.Delete(
		context.Background(),
		gophtest.AccessToken,
		id,
	)

	m.AssertExpectations(t)

	return err
}

func TestPushSecret(t *testing.T) {
	expected := uuid.NewV4()

	id, err := doPushText(t, expected, nil)

	require.NoError(t, err)
	require.Equal(t, expected, id)
}

func TestPushSecretOnRepoFailure(t *testing.T) {
	_, err := doPushText(t, uuid.UUID{}, gophtest.ErrUnexpected)

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
				{
					Id:       gophtest.CreateUUID(t, "7728154c-9400-4f1b-a2a3-01deb83ece05").String(),
					Name:     "No metadata",
					Kind:     goph.DataKind_TEXT,
					Metadata: []byte{},
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
			key := newTestKey()
			mockRV := make([]*goph.Secret, 0, len(tc.secrets))

			for _, secret := range tc.secrets {
				encrypted, err := key.Encrypt(secret.GetMetadata())
				require.NoError(t, err)

				mockRV = append(
					mockRV,
					&goph.Secret{
						Id:       secret.Id,
						Name:     secret.Name,
						Kind:     secret.Kind,
						Metadata: encrypted,
					},
				)
			}

			rv, err := doList(t, mockRV, nil)

			require.NoError(t, err)
			snaps.MatchSnapshot(t, rv)
		})
	}
}

func TestListSecretsOnDecryptFailure(t *testing.T) {
	secrets := []*goph.Secret{
		{
			Id:       gophtest.CreateUUID(t, "df566e25-43a5-4c34-9123-3931fb809b45").String(),
			Name:     gophtest.SecretName,
			Kind:     goph.DataKind_TEXT,
			Metadata: []byte(gophtest.Metadata),
		},
	}

	_, err := doList(t, secrets, nil)

	require.Error(t, err)
}

func TestListSecretsOnRepoFailure(t *testing.T) {
	_, err := doList(t, nil, gophtest.ErrUnexpected)

	require.Error(t, err)
}

func TestDeleteSecret(t *testing.T) {
	err := doDelete(t, nil)

	require.NoError(t, err)
}

func TestDeleteSecretOnRepoFailure(t *testing.T) {
	err := doDelete(t, gophtest.ErrUnexpected)

	require.Error(t, err)
}
