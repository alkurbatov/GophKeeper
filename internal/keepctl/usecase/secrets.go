package usecase

import (
	"context"
	"fmt"
	"reflect"

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

// push is low level function sending generic secret creation message to keeper.
func (uc *SecretsUseCase) push(
	ctx context.Context,
	token string,
	name string,
	kind goph.DataKind,
	description string,
	data proto.Message,
) (uuid.UUID, error) {
	var id uuid.UUID

	rawData, err := proto.Marshal(data)
	if err != nil {
		return id, fmt.Errorf("SecretsUseCase - push - proto.Marshal: %w", err)
	}

	encData, err := uc.key.Encrypt(rawData)
	if err != nil {
		return id, fmt.Errorf("SecretsUseCase - push - uc.key.Encrypt(data): %w", err)
	}

	encDescription, err := uc.key.Encrypt([]byte(description))
	if err != nil {
		return id, fmt.Errorf("SecretsUseCase - push - uc.key.Encrypt(description): %w", err)
	}

	id, err = uc.secretsRepo.Push(ctx, token, name, kind, encDescription, encData)
	if err != nil {
		return id, fmt.Errorf("SecretsUseCase - push - uc.secretsRepo.Push: %w", err)
	}

	return id, nil
}

// PushText creates new secret with arbitrary text.
func (uc *SecretsUseCase) PushText(
	ctx context.Context,
	token, name, description, text string,
) (uuid.UUID, error) {
	data := &goph.Text{
		Text: text,
	}

	return uc.push(ctx, token, name, goph.DataKind_TEXT, description, data)
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

// update is low level function sending generic secret update message to keeper.
func (uc *SecretsUseCase) update(
	ctx context.Context,
	token string,
	id uuid.UUID,
	name string,
	description string,
	noDescription bool,
	data proto.Message,
) error {
	var encData []byte

	if !reflect.ValueOf(data).IsNil() {
		rawData, err := proto.Marshal(data)
		if err != nil {
			return fmt.Errorf("SecretsUseCase - update - proto.Marshal: %w", err)
		}

		encData, err = uc.key.Encrypt(rawData)
		if err != nil {
			return fmt.Errorf("SecretsUseCase - update - uc.key.Encrypt(data): %w", err)
		}
	}

	encDescription, err := uc.key.Encrypt([]byte(description))
	if err != nil {
		return fmt.Errorf("SecretsUseCase - update - uc.key.Encrypt(description): %w", err)
	}

	if err = uc.secretsRepo.Update(
		ctx,
		token,
		id,
		name,
		encDescription,
		noDescription,
		encData,
	); err != nil {
		return fmt.Errorf("SecretsUseCase - update - uc.secretsRepo.Update: %w", err)
	}

	return nil
}

// EditText changes parameters oif stored text secret.
func (uc *SecretsUseCase) EditText(
	ctx context.Context,
	token string,
	id uuid.UUID,
	name, description string,
	noDescription bool,
	text string,
) error {
	var data *goph.Text

	if text != "" {
		data = &goph.Text{
			Text: text,
		}
	}

	return uc.update(ctx, token, id, name, description, noDescription, data)
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
