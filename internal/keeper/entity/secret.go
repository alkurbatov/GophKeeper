package entity

import (
	"errors"

	"github.com/alkurbatov/goph-keeper/pkg/goph"
	uuid "github.com/satori/go.uuid"
)

var ErrSecretExists = errors.New("secret already exists")

type Secret struct {
	ID       uuid.UUID `db:"secret_id"`
	Name     string
	Kind     goph.DataKind
	Metadata []byte
}
