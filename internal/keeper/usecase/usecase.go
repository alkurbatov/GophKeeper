package usecase

import (
	"context"

	"github.com/alkurbatov/goph-keeper/internal/keeper/config"
	"github.com/alkurbatov/goph-keeper/internal/keeper/entity"
	"github.com/alkurbatov/goph-keeper/internal/keeper/repo"
)

type Auth interface {
	Login(ctx context.Context, username, securityKey string) (entity.AccessToken, error)
}

type Secrets interface{}

type Users interface {
	Register(ctx context.Context, username, securityKey string) (entity.AccessToken, error)
}

// UseCases is a collection of business logic use cases.
type UseCases struct {
	Auth    Auth
	Secrets Secrets
	Users   Users
}

// New creates and initializes collection of business logic use cases.
func New(cfg *config.Config, repos *repo.Repositories) *UseCases {
	return &UseCases{
		Auth:    NewAuthUseCase(cfg.Secret, repos.Users),
		Secrets: NewSecretsUseCase(repos.Secrets),
		Users:   NewUsersUseCase(cfg.Secret, repos.Users),
	}
}
