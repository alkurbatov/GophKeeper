package repo

import "github.com/alkurbatov/goph-keeper/internal/keeper/infra/postgres"

var _ Users = (*UsersRepo)(nil)

// UsersRepo is facade to users stored in Postgres.
type UsersRepo struct {
	pg *postgres.Postgres
}

// NewUsersRepo creates and initializes UsersRepo object.
func NewUsersRepo(pg *postgres.Postgres) *UsersRepo {
	return &UsersRepo{pg}
}
