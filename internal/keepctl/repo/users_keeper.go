package repo

import "github.com/alkurbatov/goph-keeper/internal/keepctl/infra/grpcconn"

var _ Users = (*UsersRepo)(nil)

// UsersRepo is facade to operations regarding Keeper.
type UsersRepo struct {
	conn *grpcconn.Connection
}

// NewUsersRepo creates and initializes UsersRepo object.
func NewUsersRepo(conn *grpcconn.Connection) *UsersRepo {
	return &UsersRepo{conn}
}
