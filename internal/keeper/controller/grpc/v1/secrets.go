package v1

import (
	"context"

	"github.com/alkurbatov/goph-keeper/internal/keeper/usecase"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
)

// SecretsServer provides implementation of the Secrets API.
type SecretsServer struct {
	goph.UnimplementedSecretsServer

	secretsUseCase usecase.Secrets
}

// NewSecretsServer initializes and creates new SecretsServer.
func NewSecretsServer(secrets usecase.Secrets) *SecretsServer {
	return &SecretsServer{secretsUseCase: secrets}
}

// Create creates new secret for a user.
func (s SecretsServer) Create(
	context.Context,
	*goph.CreateSecretRequest,
) (*goph.CreateSecretResponse, error) {
	return &goph.CreateSecretResponse{}, nil
}

// List retrieves list of the secrets stored a user.
func (s SecretsServer) List(
	context.Context,
	*goph.ListSecretsRequest,
) (*goph.ListSecretsResponse, error) {
	return &goph.ListSecretsResponse{}, nil
}

// Get returns particular secret stored by a user.
func (s SecretsServer) Get(
	context.Context,
	*goph.GetSecretRequest,
) (*goph.GetSecretResponse, error) {
	return &goph.GetSecretResponse{}, nil
}

// Update updates particular secret stored by a user.
func (s SecretsServer) Update(
	context.Context,
	*goph.UpdateSecretRequest,
) (*goph.UpdateSecretResponse, error) {
	return &goph.UpdateSecretResponse{}, nil
}

// Delete remove particular secret stored by a user.
func (s SecretsServer) Delete(
	context.Context,
	*goph.DeleteSecretRequest,
) (*goph.DeleteSecretResponse, error) {
	return &goph.DeleteSecretResponse{}, nil
}
