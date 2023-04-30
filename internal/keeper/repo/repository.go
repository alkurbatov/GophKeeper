// Package repo provides facade to data stored in external sources.
package repo

import (
	"github.com/alkurbatov/goph-keeper/internal/keeper/infra/logger"
	"github.com/alkurbatov/goph-keeper/internal/keeper/infra/postgres"
)

type Secrets interface{}

type Users interface{}

// Repositories is a collection of data repositories.
type Repositories struct {
	Secrets Secrets
	Users   Users
}

// New creates and initializes collection of data repositories.
func New(log *logger.Logger, pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Secrets: NewSecretsRepo(pg),
		Users:   NewUsersRepo(pg),
	}
}
