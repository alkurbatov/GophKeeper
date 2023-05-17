package usecase

import (
	"context"
	"fmt"

	"github.com/alkurbatov/goph-keeper/internal/keeper/entity"
	"github.com/alkurbatov/goph-keeper/internal/keeper/repo"
	"github.com/alkurbatov/goph-keeper/internal/libraries/creds"
)

var _ Auth = (*AuthUseCase)(nil)

// AuthUseCase contains business logic related to authentication.
type AuthUseCase struct {
	secret    creds.Password
	usersRepo repo.Users
}

// NewAuthUseCase create and initializes new AuthUseCase object.
func NewAuthUseCase(
	secret creds.Password,
	users repo.Users,
) *AuthUseCase {
	return &AuthUseCase{secret, users}
}

func (uc *AuthUseCase) Login(
	ctx context.Context,
	username, securityKey string,
) (entity.AccessToken, error) {
	user, err := uc.usersRepo.Verify(ctx, username, securityKey)
	if err != nil {
		return "", fmt.Errorf("AuthUseCase - Login - uc.usersRepo.Verify: %w", err)
	}

	accessToken, err := entity.NewAccessToken(user, uc.secret)
	if err != nil {
		return "", fmt.Errorf("AuthUseCase - Login - entity.NewAccessToken: %w", err)
	}

	return accessToken, nil
}
