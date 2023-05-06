package repo

import (
	"context"

	"github.com/alkurbatov/goph-keeper/pkg/goph"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
)

var _ Secrets = (*SecretsRepoMock)(nil)

type SecretsRepoMock struct {
	mock.Mock
}

func (m *SecretsRepoMock) Push(
	ctx context.Context,
	token, name string,
	kind goph.DataKind,
	description, payload []byte,
) (uuid.UUID, error) {
	args := m.Called(ctx, token, name, kind, description, payload)

	return args.Get(0).(uuid.UUID), args.Error(1)
}
