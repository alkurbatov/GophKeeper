package usecase_test

import (
	"context"
	"testing"

	"github.com/alkurbatov/goph-keeper/internal/keeper/entity"
	"github.com/alkurbatov/goph-keeper/internal/keeper/repo"
	"github.com/alkurbatov/goph-keeper/internal/keeper/usecase"
	"github.com/alkurbatov/goph-keeper/internal/libraries/gophtest"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func doCreateSecret(
	t *testing.T,
	repoSecretID uuid.UUID,
	repoErr error,
) (uuid.UUID, error) {
	t.Helper()

	owner := uuid.NewV4()

	m := &repo.SecretsRepoMock{}
	m.On(
		"Create",
		mock.Anything,
		owner,
		gophtest.SecretName,
		goph.DataKind_TEXT,
		[]byte(gophtest.Metadata),
		[]byte(gophtest.TextData),
	).
		Return(repoSecretID, repoErr)

	sat := usecase.NewSecretsUseCase(m)
	id, err := sat.Create(
		context.Background(),
		owner,
		gophtest.SecretName,
		goph.DataKind_TEXT,
		[]byte(gophtest.Metadata),
		[]byte(gophtest.TextData),
	)

	m.AssertExpectations(t)

	return id, err
}

func doListSecrets(
	t *testing.T,
	repoSecrets []entity.Secret,
	repoErr error,
) ([]entity.Secret, error) {
	t.Helper()

	owner := uuid.NewV4()

	rv := make([]entity.Secret, len(repoSecrets))
	copy(rv, repoSecrets)

	m := &repo.SecretsRepoMock{}
	m.On("List", mock.Anything, owner).
		Return(rv, repoErr)

	sat := usecase.NewSecretsUseCase(m)
	secrets, err := sat.List(context.Background(), owner)

	m.AssertExpectations(t)

	return secrets, err
}

func doDeleteSecret(t *testing.T, repoErr error) error {
	t.Helper()

	owner := uuid.NewV4()
	id := uuid.NewV4()

	m := &repo.SecretsRepoMock{}
	m.On("Delete", mock.Anything, owner, id).
		Return(repoErr)

	sat := usecase.NewSecretsUseCase(m)
	err := sat.Delete(context.Background(), owner, id)

	m.AssertExpectations(t)

	return err
}

func TestCreateSecret(t *testing.T) {
	expected := uuid.NewV4()

	id, err := doCreateSecret(t, expected, nil)

	require.NoError(t, err)
	require.Equal(t, expected, id)
}

func TestCreateSecretFailsIfUserExists(t *testing.T) {
	_, err := doCreateSecret(t, uuid.UUID{}, entity.ErrSecretExists)

	require.Error(t, err)
}

func TestListSecrets(t *testing.T) {
	tt := []struct {
		name     string
		expected []entity.Secret
	}{
		{
			name: "List secrets of a user",
			expected: []entity.Secret{
				{ID: uuid.NewV4(), Name: gophtest.SecretName, Kind: goph.DataKind_BINARY},
				{
					ID:       uuid.NewV4(),
					Name:     gophtest.SecretName + "ex",
					Kind:     goph.DataKind_TEXT,
					Metadata: []byte(gophtest.Metadata),
				},
			},
		},
		{
			name:     "List secrets when user has no secrets",
			expected: []entity.Secret{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rv, err := doListSecrets(t, tc.expected, nil)

			require.NoError(t, err)
			require.Equal(t, tc.expected, rv)
		})
	}
}

func TestListSecretFailsOnRepoFailure(t *testing.T) {
	_, err := doListSecrets(t, nil, gophtest.ErrUnexpected)

	require.Error(t, err)
}

func TestDeleteSecret(t *testing.T) {
	err := doDeleteSecret(t, nil)

	require.NoError(t, err)
}

func TestDeleteSecretFailsOnRepoFailure(t *testing.T) {
	err := doDeleteSecret(t, gophtest.ErrUnexpected)

	require.Error(t, err)
}
