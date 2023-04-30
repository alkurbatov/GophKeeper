package usecase

import "github.com/alkurbatov/goph-keeper/internal/keepctl/repo"

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
