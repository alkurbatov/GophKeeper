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

func (m *SecretsRepoMock) Create(
	ctx context.Context,
	owner uuid.UUID,
	name string,
	kind goph.DataKind,
	metadata, data []byte,
) (uuid.UUID, error) {
	args := m.Called(ctx, owner, name, kind, metadata, data)

	return args.Get(0).(uuid.UUID), args.Error(1)
}
