package usecase

import "github.com/alkurbatov/goph-keeper/internal/keepctl/repo"

var _ Users = (*UsersUseCase)(nil)

// UsersUseCase contains business logic related to users management.
type UsersUseCase struct {
	usersRepo repo.Users
}

// NewUsersUseCase create and initializes new UsersUseCase object.
func NewUsersUseCase(users repo.Users) *UsersUseCase {
	return &UsersUseCase{users}
}
