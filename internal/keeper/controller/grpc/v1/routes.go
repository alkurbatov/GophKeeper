// Package v1 implements v1 version of the gRPC API.
package v1

import (
	"github.com/alkurbatov/goph-keeper/internal/keeper/usecase"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
	"google.golang.org/grpc"
)

// RegisterRoutes injects new routes into the provided gRPC server.
func RegisterRoutes(server *grpc.Server, useCases *usecase.UseCases) {
	auth := NewAuthServer(useCases.Auth)
	goph.RegisterAuthServer(server, auth)

	secrets := NewSecretsServer(useCases.Secrets)
	goph.RegisterSecretsServer(server, secrets)

	users := NewUsersServer(useCases.Users)
	goph.RegisterUsersServer(server, users)
}
