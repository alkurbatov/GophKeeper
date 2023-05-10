package entity

import (
	"errors"

	"github.com/alkurbatov/goph-keeper/pkg/goph"
	uuid "github.com/satori/go.uuid"
)

var (
	ErrSecretNotFound = errors.New("secret not found")
	ErrSecretExists   = errors.New("secret already exists")
)

// Secret represents full secret info stored in the service.
type Secret struct {
	ID       uuid.UUID `db:"secret_id"`
	Name     string
	Kind     goph.DataKind
	Metadata []byte
	Data     []byte
}
