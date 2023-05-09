package v1

import (
	"context"
	"errors"

	"github.com/alkurbatov/goph-keeper/internal/keeper/entity"
	"github.com/alkurbatov/goph-keeper/internal/keeper/usecase"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	ctx context.Context,
	req *goph.CreateSecretRequest,
) (*goph.CreateSecretResponse, error) {
	name := req.GetName()
	kind := req.GetKind()
	metadata := req.GetMetadata()
	data := req.GetData()

	if details, ok := validateSecret(name, metadata, data); !ok {
		st := composeBadRequestError(details)

		return nil, st.Err()
	}

	owner := entity.UserFromContext(ctx)
	if owner == nil {
		return nil, status.Errorf(codes.Unauthenticated, entity.ErrInvalidCredentials.Error())
	}

	id, err := s.secretsUseCase.Create(ctx, owner.ID, name, kind, metadata, data)
	if err != nil {
		if errors.Is(err, entity.ErrSecretExists) {
			return nil, status.Errorf(codes.AlreadyExists, entity.ErrSecretExists.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &goph.CreateSecretResponse{Id: id.String()}, nil
}

// List retrieves list of the secrets stored a user.
func (s SecretsServer) List(
	ctx context.Context,
	_ *goph.ListSecretsRequest,
) (*goph.ListSecretsResponse, error) {
	owner := entity.UserFromContext(ctx)
	if owner == nil {
		return nil, status.Errorf(codes.Unauthenticated, entity.ErrInvalidCredentials.Error())
	}

	data, err := s.secretsUseCase.List(ctx, owner.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	rv := make([]*goph.Secret, 0, len(data))
	for _, val := range data {
		rv = append(rv, &goph.Secret{
			Id:       val.ID.String(),
			Name:     val.Name,
			Kind:     val.Kind,
			Metadata: val.Metadata,
		})
	}

	return &goph.ListSecretsResponse{Secrets: rv}, nil
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
