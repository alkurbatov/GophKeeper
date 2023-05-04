package repo

import "github.com/alkurbatov/goph-keeper/pkg/goph"

var _ Secrets = (*SecretsRepo)(nil)

// SecretsRepo is facade to secrets stored in Keeper.
type SecretsRepo struct {
	client goph.SecretsClient
}

// NewSecretsRepo creates and initializes SecretsRepo object.
func NewSecretsRepo(client goph.SecretsClient) *SecretsRepo {
	return &SecretsRepo{client}
}
