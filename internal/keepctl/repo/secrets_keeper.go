package repo

import "github.com/alkurbatov/goph-keeper/internal/keepctl/infra/grpcconn"

var _ Secrets = (*SecretsRepo)(nil)

// SecretsRepo is facade to secrets stored in Keeper.
type SecretsRepo struct {
	conn *grpcconn.Connection
}

// NewSecretsRepo creates and initializes SecretsRepo object.
func NewSecretsRepo(conn *grpcconn.Connection) *SecretsRepo {
	return &SecretsRepo{conn}
}
