package usecase

import (
	"context"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/repo"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/proto"
)

type Auth interface {
	Login(ctx context.Context, username string, key entity.Key) (string, error)
}

type Secrets interface {
	PushBinary(ctx context.Context, token, name, description string, binary []byte) (uuid.UUID, error)
	PushCreds(ctx context.Context, token, name, description, login, password string) (uuid.UUID, error)
	PushText(ctx context.Context, token, name, description, text string) (uuid.UUID, error)

	List(ctx context.Context, token string) ([]*goph.Secret, error)
	Get(ctx context.Context, token string, id uuid.UUID) (*goph.Secret, proto.Message, error)

	EditBinary(
		ctx context.Context,
		token string,
		id uuid.UUID,
		name, description string,
		noDescription bool,
		binary []byte,
	) error

	EditCreds(
		ctx context.Context,
		token string,
		id uuid.UUID,
		name, description string,
		noDescription bool,
		login, password string,
	) error

	EditText(
		ctx context.Context,
		token string,
		id uuid.UUID,
		name, description string,
		noDescription bool,
		text string,
	) error

	Delete(ctx context.Context, token string, id uuid.UUID) error
}

type Users interface {
	Register(ctx context.Context, username string, key entity.Key) (string, error)
}

// UseCases is a collection of business logic use cases.
type UseCases struct {
	Auth    Auth
	Secrets Secrets
	Users   Users
}

// New creates and initializes collection of business logic use cases.
func New(key entity.Key, repos *repo.Repositories) *UseCases {
	return &UseCases{
		Auth:    NewAuthUseCase(repos.Auth),
		Secrets: NewSecretsUseCase(key, repos.Secrets),
		Users:   NewUsersUseCase(repos.Users),
	}
}
