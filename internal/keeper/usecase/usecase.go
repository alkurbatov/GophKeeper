package usecase

import (
	"github.com/alkurbatov/goph-keeper/internal/keeper/config"
	"github.com/alkurbatov/goph-keeper/internal/keeper/infra/logger"
	"github.com/alkurbatov/goph-keeper/internal/keeper/repo"
)

type Auth interface{}

type Secrets interface{}

type Users interface{}

// UseCases is a collection of business logic use cases.
type UseCases struct {
	Auth    Auth
	Secrets Secrets
	Users   Users
}

// New creates and initializes collection of business logic use cases.
func New(log *logger.Logger, cfg *config.Config, repos *repo.Repositories) *UseCases {
	return &UseCases{
		Auth:    NewAuthUseCase(log, cfg.Secret, repos.Users),
		Secrets: NewSecretsUseCase(repos.Secrets),
		Users:   NewUsersUseCase(repos.Users),
	}
}
