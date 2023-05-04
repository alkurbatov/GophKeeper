package usecase

import (
	"context"
	"fmt"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/repo"
)

var _ Auth = (*AuthUseCase)(nil)

// AuthUseCase contains business logic related to authentication.
type AuthUseCase struct {
	authRepo repo.Auth
}

// NewAuthUseCase create and initializes new AuthUseCase object.
func NewAuthUseCase(
	auth repo.Auth,
) *AuthUseCase {
	return &AuthUseCase{auth}
}

// Login authenticates a user.
func (uc *AuthUseCase) Login(
	ctx context.Context,
	username, password string,
) (string, error) {
	key := entity.NewKey(username, password)
	securityKey := key.Hash()

	token, err := uc.authRepo.Login(ctx, username, securityKey)
	if err != nil {
		return "", fmt.Errorf("AuthUseCase - Login - uc.authRepo.Login: %w", err)
	}

	return token, nil
}
