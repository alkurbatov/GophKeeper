package repo

import "github.com/alkurbatov/goph-keeper/internal/keeper/infra/postgres"

var _ Secrets = (*SecretsRepo)(nil)

// SecretsRepo is facade to secrets stored in Postgres.
type SecretsRepo struct {
	pg *postgres.Postgres
}

// NewSecretsRepo creates and initializes SecretsRepo object.
func NewSecretsRepo(pg *postgres.Postgres) *SecretsRepo {
	return &SecretsRepo{pg}
}
