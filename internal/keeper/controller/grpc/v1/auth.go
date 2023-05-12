package v1

import (
	"context"

	"github.com/alkurbatov/goph-keeper/internal/keeper/usecase"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
)

// AuthServer provides implementation of the Auth API.
type AuthServer struct {
	goph.UnimplementedAuthServer

	authUseCase usecase.Auth
}

// NewAuthServer initializes and creates new AuthServer.
func NewAuthServer(auth usecase.Auth) *AuthServer {
	return &AuthServer{authUseCase: auth}
}

// Login authenticates a user into the service.
func (s AuthServer) Login(
	context.Context,
	*goph.LoginRequest,
) (*goph.LoginResponse, error) {
	return &goph.LoginResponse{}, nil
}
