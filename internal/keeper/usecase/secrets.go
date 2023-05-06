package usecase

import (
	"context"
	"fmt"

	"github.com/alkurbatov/goph-keeper/internal/keeper/repo"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
	uuid "github.com/satori/go.uuid"
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

func (uc *SecretsUseCase) Create(
	ctx context.Context,
	owner uuid.UUID,
	name string,
	kind goph.DataKind,
	metadata, data []byte,
) (uuid.UUID, error) {
	id, err := uc.secretsRepo.Create(ctx, owner, name, kind, metadata, data)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("SecretsUseCase - Create - uc.secretsRepo.Create: %w", err)
	}

	return id, nil
}
