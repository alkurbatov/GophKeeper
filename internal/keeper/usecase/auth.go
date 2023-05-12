package usecase

import (
	"github.com/alkurbatov/goph-keeper/internal/keeper/entity"
	"github.com/alkurbatov/goph-keeper/internal/keeper/infra/logger"
	"github.com/alkurbatov/goph-keeper/internal/keeper/repo"
)

var _ Auth = (*AuthUseCase)(nil)

const _defaultMinimalSecretLength = 32

// AuthUseCase contains business logic related to authentication.
type AuthUseCase struct {
	secret    entity.Secret
	usersRepo repo.Users
}

// NewAuthUseCase create and initializes new AuthUseCase object.
func NewAuthUseCase(
	log *logger.Logger,
	secret entity.Secret,
	users repo.Users,
) *AuthUseCase {
	if len([]byte(secret)) < _defaultMinimalSecretLength {
		log.Warn().Msg("Insecure signature: secret key is shorter than 32 bytes!")
	}

	return &AuthUseCase{secret, users}
}
