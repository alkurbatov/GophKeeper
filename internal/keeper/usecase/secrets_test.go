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
