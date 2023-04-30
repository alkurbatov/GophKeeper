package usecase

import "github.com/alkurbatov/goph-keeper/internal/keepctl/repo"

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
func New(repos *repo.Repositories) *UseCases {
	return &UseCases{
		Auth:    NewAuthUseCase(repos.Users),
		Secrets: NewSecretsUseCase(repos.Secrets),
		Users:   NewUsersUseCase(repos.Users),
	}
}
