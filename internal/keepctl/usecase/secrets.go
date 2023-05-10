package usecase

import (
	"context"
	"fmt"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/repo"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/proto"
)

var _ Secrets = (*SecretsUseCase)(nil)

// SecretsUseCase contains business logic related to secrets management.
type SecretsUseCase struct {
	key         entity.Key
	secretsRepo repo.Secrets
}

// NewSecretsUseCase create and initializes new SecretsUseCase object.
func NewSecretsUseCase(key entity.Key, secrets repo.Secrets) *SecretsUseCase {
	return &SecretsUseCase{key, secrets}
}

// push is low level push function sending generic message to keeper.
func (uc *SecretsUseCase) push(
	ctx context.Context,
	token string,
	name string,
	data proto.Message,
	description string,
) (uuid.UUID, error) {
	var id uuid.UUID

	rawData, err := proto.Marshal(data)
	if err != nil {
		return id, fmt.Errorf("SecretsUseCase - PushText - proto.Marshal: %w", err)
	}

	payload, err := uc.key.Encrypt(rawData)
	if err != nil {
		return id, fmt.Errorf("SecretsUseCase - PushText - uc.key.Encrypt: %w", err)
	}

	metadata, err := uc.key.Encrypt([]byte(description))
	if err != nil {
		return id, fmt.Errorf("SecretsUseCase - PushText - uc.key.Encrypt: %w", err)
	}

	id, err = uc.secretsRepo.Push(ctx, token, name, goph.DataKind_TEXT, metadata, payload)
	if err != nil {
		return id, fmt.Errorf("SecretsUseCase - PushText - uc.secretsRepo.Push: %w", err)
	}

	return id, nil
}

// PushText creates new secret with arbitrary text.
func (uc *SecretsUseCase) PushText(
	ctx context.Context,
	token, name, text, description string,
) (uuid.UUID, error) {
	data := &goph.Text{
		Text: text,
	}

	return uc.push(ctx, token, name, data, description)
}

// List returns list of user's secrets.
// All sensitive parts are decrypted.
func (uc *SecretsUseCase) List(ctx context.Context, token string) ([]*goph.Secret, error) {
	data, err := uc.secretsRepo.List(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("SecretsUseCase - List - uc.secretsRepo.List: %w", err)
	}

	for i, val := range data {
		data[i].Metadata, err = uc.key.Decrypt(val.GetMetadata())
		if err != nil {
			return nil, fmt.Errorf("SecretsUseCase - List - uc.key.Decrypt: %w", err)
		}
	}

	return data, nil
}

// Get retrieves full user's secret.
// All sensitive parts are decrypted.
func (uc *SecretsUseCase) Get(
	ctx context.Context,
	token string,
	id uuid.UUID,
) (*goph.Secret, []byte, error) {
	secret, data, err := uc.secretsRepo.Get(ctx, token, id)
	if err != nil {
		return nil, nil, fmt.Errorf("SecretsUseCase - Get - uc.secretsRepo.Get: %w", err)
	}

	secret.Metadata, err = uc.key.Decrypt(secret.GetMetadata())
	if err != nil {
		return nil, nil, fmt.Errorf("SecretsUseCase - Get - uc.key.Decrypt(metadata): %w", err)
	}

	decryptedData, err := uc.key.Decrypt(data)
	if err != nil {
		return nil, nil, fmt.Errorf("SecretsUseCase - Get - uc.key.Decrypt(data): %w", err)
	}

	return secret, decryptedData, nil
}

// Delete removes user's secret.
func (uc *SecretsUseCase) Delete(
	ctx context.Context,
	token string,
	id uuid.UUID,
) error {
	if err := uc.secretsRepo.Delete(ctx, token, id); err != nil {
		return fmt.Errorf("SecretsUseCase - Delete - uc.secretsRepo.Delete: %w", err)
	}

	return nil
}
