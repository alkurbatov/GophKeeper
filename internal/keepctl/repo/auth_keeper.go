package repo

import "github.com/alkurbatov/goph-keeper/internal/keepctl/infra/grpcconn"

var _ Auth = (*AuthRepo)(nil)

// AuthRepo is facade to operations regarding Keeper.
type AuthRepo struct {
	conn *grpcconn.Connection
}

// NewAuthRepo creates and initializes AuthRepo object.
func NewAuthRepo(conn *grpcconn.Connection) *AuthRepo {
	return &AuthRepo{conn}
}
