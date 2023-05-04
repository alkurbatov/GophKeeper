package usecase

import (
	"context"
	"fmt"

	"github.com/alkurbatov/goph-keeper/internal/keeper/entity"
	"github.com/alkurbatov/goph-keeper/internal/keeper/repo"
)

var _ Users = (*UsersUseCase)(nil)

// UsersUseCase contains business logic related to users management.
type UsersUseCase struct {
	secret    entity.Secret
	usersRepo repo.Users
}

// NewUsersUseCase create and initializes new UsersUseCase object.
func NewUsersUseCase(secret entity.Secret, users repo.Users) *UsersUseCase {
	return &UsersUseCase{secret, users}
}

// Register creates a new user.
func (uc UsersUseCase) Register(
	ctx context.Context,
	username, securityKey string,
) (entity.AccessToken, error) {
	id, err := uc.usersRepo.Register(ctx, username, securityKey)
	if err != nil {
		return "", fmt.Errorf("UsersUseCase - Register - uc.usersRepo.Register: %w", err)
	}

	user := entity.User{
		ID:       id,
		Username: username,
	}

	accessToken, err := entity.NewAccessToken(user, uc.secret)
	if err != nil {
		return "", fmt.Errorf("AuthUseCase - Login - entity.NewAccessToken: %w", err)
	}

	return accessToken, nil
}
