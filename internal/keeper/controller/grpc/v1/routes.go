// Package v1 implements v1 version of the gRPC API.
package v1

import (
	"github.com/alkurbatov/goph-keeper/internal/keeper/infra/grpcserver"
	"github.com/alkurbatov/goph-keeper/internal/keeper/usecase"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
)

// RegisterRoutes injects new routes into the provided gRPC server.
func RegisterRoutes(server *grpcserver.Server, useCases *usecase.UseCases) {
	auth := NewAuthServer(useCases.Auth)
	goph.RegisterAuthServer(server.Instance(), auth)

	secrets := NewSecretsServer(useCases.Secrets)
	goph.RegisterSecretsServer(server.Instance(), secrets)

	users := NewUsersServer(useCases.Users)
	goph.RegisterUsersServer(server.Instance(), users)
}
