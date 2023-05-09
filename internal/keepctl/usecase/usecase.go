package usecase

import (
	"context"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/repo"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
	uuid "github.com/satori/go.uuid"
)

type Auth interface {
	Login(ctx context.Context, username string, key entity.Key) (string, error)
}

type Secrets interface {
	PushText(ctx context.Context, token, name, text, description string) (uuid.UUID, error)
	List(ctx context.Context, token string) ([]*goph.Secret, error)
	Delete(ctx context.Context, token string, id uuid.UUID) error
}

type Users interface {
	Register(ctx context.Context, username string, key entity.Key) (string, error)
}

// UseCases is a collection of business logic use cases.
type UseCases struct {
	Auth    Auth
	Secrets Secrets
	Users   Users
}

// New creates and initializes collection of business logic use cases.
func New(key entity.Key, repos *repo.Repositories) *UseCases {
	return &UseCases{
		Auth:    NewAuthUseCase(repos.Auth),
		Secrets: NewSecretsUseCase(key, repos.Secrets),
		Users:   NewUsersUseCase(repos.Users),
	}
}
