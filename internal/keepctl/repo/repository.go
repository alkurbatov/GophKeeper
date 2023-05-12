package repo

import (
	"github.com/alkurbatov/goph-keeper/internal/keepctl/infra/grpcconn"
)

type Auth interface{}

type Secrets interface{}

type Users interface{}

// Repositories is a collection of data repositories.
type Repositories struct {
	Auth    Auth
	Secrets Secrets
	Users   Users
}

// New creates and initializes collection of data repositories.
func New(conn *grpcconn.Connection) *Repositories {
	return &Repositories{
		Auth:    NewAuthRepo(conn),
		Secrets: NewSecretsRepo(conn),
		Users:   NewUsersRepo(conn),
	}
}
