package usecase

import (
	"context"

	"github.com/alkurbatov/goph-keeper/internal/keeper/entity"
	"github.com/stretchr/testify/mock"
)

var _ Users = (*UsersUseCaseMock)(nil)

type UsersUseCaseMock struct {
	mock.Mock
}

func (m *UsersUseCaseMock) Register(
	ctx context.Context,
	username, securityKey string,
) (entity.AccessToken, error) {
	args := m.Called(ctx, username, securityKey)

	return args.Get(0).(entity.AccessToken), args.Error(1)
}
