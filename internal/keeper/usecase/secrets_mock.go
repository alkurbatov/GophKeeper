package usecase

import (
	"context"

	"github.com/alkurbatov/goph-keeper/internal/keeper/entity"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
)

var _ Secrets = (*SecretsUseCaseMock)(nil)

type SecretsUseCaseMock struct {
	mock.Mock
}

func (m *SecretsUseCaseMock) Create(
	ctx context.Context,
	owner uuid.UUID,
	name string,
	kind goph.DataKind,
	metadata, data []byte,
) (uuid.UUID, error) {
	args := m.Called(ctx, owner, name, kind, metadata, data)

	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *SecretsUseCaseMock) List(
	ctx context.Context,
	owner uuid.UUID,
) ([]entity.Secret, error) {
	args := m.Called(ctx, owner)

	if args.Get(0) == 1 {
		return nil, args.Error(1)
	}

	return args.Get(0).([]entity.Secret), args.Error(1)
}

func (m *SecretsUseCaseMock) Get(
	ctx context.Context,
	owner, id uuid.UUID,
) (*entity.Secret, error) {
	args := m.Called(ctx, owner, id)

	return args.Get(0).(*entity.Secret), args.Error(1)
}

func (m *SecretsUseCaseMock) Update(
	ctx context.Context,
	owner, id uuid.UUID,
	changed []string,
	name string,
	metadata []byte,
	data []byte,
) error {
	args := m.Called(ctx, owner, id, changed, name, metadata, data)

	return args.Error(0)
}

func (m *SecretsUseCaseMock) Delete(
	ctx context.Context,
	owner, id uuid.UUID,
) error {
	args := m.Called(ctx, owner, id)

	return args.Error(0)
}
