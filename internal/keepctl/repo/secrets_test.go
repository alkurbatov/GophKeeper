package repo_test

import (
	"context"
	"testing"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/repo"
	"github.com/alkurbatov/goph-keeper/internal/libraries/gophtest"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newCreateSecretRequest() *goph.CreateSecretRequest {
	return &goph.CreateSecretRequest{
		Name:     gophtest.SecretName,
		Metadata: []byte(gophtest.Metadata),
		Kind:     goph.DataKind_TEXT,
		Data:     []byte(gophtest.TextData),
	}
}

func TestCreate(t *testing.T) {
	expected := uuid.NewV4()

	resp := &goph.CreateSecretResponse{
		Id: expected.String(),
	}

	m := &goph.SecretsClientMock{}
	m.On(
		"Create",
		mock.Anything,
		newCreateSecretRequest(),
		mock.Anything,
	).
		Return(resp, nil)

	sat := repo.NewSecretsRepo(m)
	id, err := sat.Push(
		context.Background(),
		gophtest.AccessToken,
		gophtest.SecretName,
		goph.DataKind_TEXT,
		[]byte(gophtest.Metadata),
		[]byte(gophtest.TextData),
	)

	require.NoError(t, err)
	require.Equal(t, expected, id)
	m.AssertExpectations(t)
}

func TestCreateOnOperationFailure(t *testing.T) {
	m := &goph.SecretsClientMock{}
	m.On(
		"Create",
		mock.Anything,
		newCreateSecretRequest(),
		mock.Anything,
	).
		Return(nil, gophtest.ErrUnexpected)

	sat := repo.NewSecretsRepo(m)
	_, err := sat.Push(
		context.Background(),
		gophtest.AccessToken,
		gophtest.SecretName,
		goph.DataKind_TEXT,
		[]byte(gophtest.Metadata),
		[]byte(gophtest.TextData),
	)

	require.Error(t, err)
	m.AssertExpectations(t)
}
