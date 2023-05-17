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

func doGetSecret(t *testing.T, repoRV *entity.Secret, repoErr error) (*entity.Secret, error) {
	t.Helper()

	owner := uuid.NewV4()
	id := uuid.NewV4()

	m := &repo.SecretsRepoMock{}
	m.On("Get", mock.Anything, owner, id).
		Return(repoRV, repoErr)

	sat := usecase.NewSecretsUseCase(m)
	secret, err := sat.Get(context.Background(), owner, id)

	m.AssertExpectations(t)

	return secret, err
}

func doUpdateSecret(t *testing.T, repoErr error) error {
	t.Helper()

	owner := uuid.NewV4()
	id := uuid.NewV4()
	changed := []string{"name", "metadata", "data"}

	m := &repo.SecretsRepoMock{}
	m.On(
		"Update",
		mock.Anything,
		owner,
		id,
		changed,
		gophtest.SecretName,
		[]byte(gophtest.Metadata),
		[]byte(gophtest.TextData),
	).
		Return(repoErr)

	sat := usecase.NewSecretsUseCase(m)
	err := sat.Update(
		context.Background(),
		owner,
		id,
		changed,
		gophtest.SecretName,
		[]byte(gophtest.Metadata),
		[]byte(gophtest.TextData),
	)

	m.AssertExpectations(t)

	return err
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
	type expected struct {
		id  uuid.UUID
		err error
	}

	tt := []struct {
		name     string
		expected expected
	}{
		{
			name: "Create secret",
			expected: expected{
				id:  uuid.NewV4(),
				err: nil,
			},
		},
		{
			name: "Create secret fails if secret exists",
			expected: expected{
				id:  uuid.UUID{},
				err: entity.ErrSecretExists,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			id, err := doCreateSecret(t, tc.expected.id, tc.expected.err)

			require.ErrorIs(t, err, tc.expected.err)
			require.Equal(t, tc.expected.id, id)
		})
	}
}

func TestListSecrets(t *testing.T) {
	type expected struct {
		secrets []entity.Secret
		err     error
	}

	tt := []struct {
		name     string
		expected expected
	}{
		{
			name: "List secrets of a user",
			expected: expected{
				secrets: []entity.Secret{
					{
						ID:   uuid.NewV4(),
						Name: gophtest.SecretName,
						Kind: goph.DataKind_BINARY,
					},
					{
						ID:       uuid.NewV4(),
						Name:     gophtest.SecretName + "ex",
						Kind:     goph.DataKind_TEXT,
						Metadata: []byte(gophtest.Metadata),
					},
				},
				err: nil,
			},
		},
		{
			name: "List secrets when user has no secrets",
			expected: expected{
				secrets: []entity.Secret{},
				err:     nil,
			},
		},
		{
			name: "List secrets fails if repo fails",
			expected: expected{
				secrets: nil,
				err:     gophtest.ErrUnexpected,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rv, err := doListSecrets(t, tc.expected.secrets, tc.expected.err)

			require.ErrorIs(t, err, tc.expected.err)
			require.Equal(t, tc.expected.secrets, rv)
		})
	}
}

func TestGetSecret(t *testing.T) {
	type expected struct {
		secret *entity.Secret
		err    error
	}

	tt := []struct {
		name     string
		expected expected
	}{
		{
			name: "Get secret",
			expected: expected{
				secret: &entity.Secret{
					ID:       uuid.NewV4(),
					Name:     gophtest.SecretName,
					Kind:     goph.DataKind_TEXT,
					Metadata: []byte(gophtest.Metadata),
					Data:     []byte(gophtest.TextData),
				},
				err: nil,
			},
		},
		{
			name: "Get secret fails if secret doesn't exist",
			expected: expected{
				secret: nil,
				err:    entity.ErrSecretNotFound,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			secret, err := doGetSecret(t, tc.expected.secret, tc.expected.err)

			require.ErrorIs(t, err, tc.expected.err)
			require.Equal(t, tc.expected.secret, secret)
		})
	}
}

func TestUpdateSecret(t *testing.T) {
	tt := []struct {
		name     string
		expected error
	}{
		{
			name:     "Update secret",
			expected: nil,
		},
		{
			name:     "Update secret if secret not found",
			expected: entity.ErrSecretNotFound,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := doUpdateSecret(t, tc.expected)

			require.ErrorIs(t, err, tc.expected)
		})
	}
}

func TestDeleteSecret(t *testing.T) {
	tt := []struct {
		name     string
		expected error
	}{
		{
			name:     "Delete secret",
			expected: nil,
		},
		{
			name:     "Delete secret if secret not found",
			expected: entity.ErrSecretNotFound,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := doDeleteSecret(t, tc.expected)

			require.ErrorIs(t, err, tc.expected)
		})
	}
}
