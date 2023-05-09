package repo

import (
	"context"
	"fmt"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/metadata"
)

var _ Secrets = (*SecretsRepo)(nil)

// SecretsRepo is facade to secrets stored in Keeper.
type SecretsRepo struct {
	client goph.SecretsClient
}

// NewSecretsRepo creates and initializes SecretsRepo object.
func NewSecretsRepo(client goph.SecretsClient) *SecretsRepo {
	return &SecretsRepo{client}
}

// Push send new secret data to the server.
func (r *SecretsRepo) Push(
	ctx context.Context,
	token, name string,
	kind goph.DataKind,
	description, payload []byte,
) (uuid.UUID, error) {
	var id uuid.UUID

	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	req := &goph.CreateSecretRequest{
		Name:     name,
		Metadata: description,
		Kind:     kind,
		Data:     payload,
	}

	resp, err := r.client.Create(ctx, req)
	if err != nil {
		return id, fmt.Errorf("SecretsRepo - Push - r.client.Create: %w", entity.NewRequestError(err))
	}

	id, err = uuid.FromString(resp.GetId())
	if err != nil {
		return id, fmt.Errorf("SecretsRepo - Push - uuid.FromString: %w", err)
	}

	return id, nil
}

// List returns list of  user's secrets.
func (r *SecretsRepo) List(
	ctx context.Context,
	token string,
) ([]*goph.Secret, error) {
	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	req := &goph.ListSecretsRequest{}

	resp, err := r.client.List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("SecretsRepo - List - r.client.List: %w", entity.NewRequestError(err))
	}

	return resp.Secrets, nil
}
