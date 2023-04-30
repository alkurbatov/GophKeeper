package v1

import (
	"context"

	"github.com/alkurbatov/goph-keeper/internal/keeper/usecase"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
)

// UsersServer provides implementation of the Users API.
type UsersServer struct {
	goph.UnimplementedUsersServer

	usersUseCase usecase.Users
}

// NewUsersServer initializes and creates new UsersServer.
func NewUsersServer(users usecase.Users) *UsersServer {
	return &UsersServer{usersUseCase: users}
}

// Register creates new user.
func (s UsersServer) Register(
	context.Context,
	*goph.RegisterUserRequest,
) (*goph.RegisterUserResponse, error) {
	return &goph.RegisterUserResponse{}, nil
}
