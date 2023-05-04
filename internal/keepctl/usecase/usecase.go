package usecase

import (
	"context"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/repo"
)

type Auth interface {
	Login(ctx context.Context, username, password string) (string, error)
}

type Secrets interface{}

type Users interface {
	Register(ctx context.Context, username, password string) (string, error)
}

// UseCases is a collection of business logic use cases.
type UseCases struct {
	Auth    Auth
	Secrets Secrets
	Users   Users
}

// New creates and initializes collection of business logic use cases.
func New(repos *repo.Repositories) *UseCases {
	return &UseCases{
		Auth:    NewAuthUseCase(repos.Auth),
		Secrets: NewSecretsUseCase(repos.Secrets),
		Users:   NewUsersUseCase(repos.Users),
	}
}
