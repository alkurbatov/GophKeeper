package usecase

import (
	"github.com/alkurbatov/goph-keeper/internal/keeper/repo"
)

var _ Secrets = (*SecretsUseCase)(nil)

// SecretsUseCase contains business logic related to secrets management.
type SecretsUseCase struct {
	secretsRepo repo.Secrets
}

// NewSecretsUseCase create and initializes new SecretsUseCase object.
func NewSecretsUseCase(secrets repo.Secrets) *SecretsUseCase {
	return &SecretsUseCase{secrets}
}
