package repo

import (
	"context"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/infra/grpcconn"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
)

type Auth interface {
	Login(ctx context.Context, username, securityKey string) (string, error)
}

type Secrets interface{}

type Users interface {
	Register(ctx context.Context, username, securityKey string) (string, error)
}

// Repositories is a collection of data repositories.
type Repositories struct {
	Auth    Auth
	Secrets Secrets
	Users   Users
}

// New creates and initializes collection of data repositories.
func New(conn *grpcconn.Connection) *Repositories {
	c := conn.Instance()

	return &Repositories{
		Auth:    NewAuthRepo(goph.NewAuthClient(c)),
		Secrets: NewSecretsRepo(goph.NewSecretsClient(c)),
		Users:   NewUsersRepo(goph.NewUsersClient(c)),
	}
}
