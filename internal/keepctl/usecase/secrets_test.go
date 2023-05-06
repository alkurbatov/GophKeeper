package usecase_test

import (
	"context"
	"testing"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/repo"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/usecase"
	"github.com/alkurbatov/goph-keeper/internal/libraries/gophtest"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPush(t *testing.T) {
	expected := uuid.NewV4()

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
		Return(expected, nil)

	sat := usecase.NewSecretsUseCase(newTestKey(), m)
	id, err := sat.PushText(
		context.Background(),
		gophtest.AccessToken,
		gophtest.SecretName,
		gophtest.Metadata,
		gophtest.TextData,
	)

	require.NoError(t, err)
	require.Equal(t, expected, id)
	m.AssertExpectations(t)
}

func TestPushOnRepoFailure(t *testing.T) {
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
		Return(uuid.UUID{}, gophtest.ErrUnexpected)

	sat := usecase.NewSecretsUseCase(newTestKey(), m)
	_, err := sat.PushText(
		context.Background(),
		gophtest.AccessToken,
		gophtest.SecretName,
		gophtest.Metadata,
		gophtest.TextData,
	)

	require.Error(t, err)
	m.AssertExpectations(t)
}
