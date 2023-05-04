package entity

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or security key")
	ErrNoSecurityKey      = errors.New("security key not set")
	ErrUserExists         = errors.New("user already exists")
)

// User represents basic user of the system.
type User struct {
	ID       uuid.UUID `db:"user_id"`
	Username string
}
