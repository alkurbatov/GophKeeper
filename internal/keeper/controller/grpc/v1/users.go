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
	ctx context.Context,
	req *goph.RegisterUserRequest,
) (*goph.RegisterUserResponse, error) {
	username := req.GetUsername()
	key := req.GetSecurityKey()

	if details, ok := validateCredentials(username, key); !ok {
		st := composeBadRequestError(details)

		return nil, st.Err()
	}

	accessToken, err := s.usersUseCase.Register(ctx, username, key)
	if err != nil {
		if errors.Is(err, entity.ErrUserExists) {
			return nil, status.Errorf(codes.AlreadyExists, entity.ErrUserExists.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &goph.RegisterUserResponse{AccessToken: accessToken.String()}, nil
}
