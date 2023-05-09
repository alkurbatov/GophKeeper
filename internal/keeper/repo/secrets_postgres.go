package repo

import (
	"context"
	"fmt"

	"github.com/alkurbatov/goph-keeper/internal/keeper/entity"
	"github.com/alkurbatov/goph-keeper/internal/keeper/infra/logger"
	"github.com/alkurbatov/goph-keeper/internal/keeper/infra/postgres"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
	uuid "github.com/satori/go.uuid"
)

var _ Secrets = (*SecretsRepo)(nil)

// SecretsRepo is facade to secrets stored in Postgres.
type SecretsRepo struct {
	pg *postgres.Postgres
}

// NewSecretsRepo creates and initializes SecretsRepo object.
func NewSecretsRepo(pg *postgres.Postgres) *SecretsRepo {
	return &SecretsRepo{pg}
}

// Create stores new secret in database.
func (r *SecretsRepo) Create(
	ctx context.Context,
	owner uuid.UUID,
	name string,
	kind goph.DataKind,
	metadata, data []byte,
) (uuid.UUID, error) {
	var id uuid.UUID

	tx, err := r.pg.BeginTx(ctx)
	if err != nil {
		return id, fmt.Errorf("SecretsRepo - Create - r.pg.BeginTx: %w", err)
	}

	defer func() {
		switch err {
		case nil:
			if cErr := tx.Commit(context.Background()); cErr != nil {
				err = fmt.Errorf("SecretsRepo - Create - tx.Commit: %w", cErr)
			}
		default:
			if rErr := tx.Rollback(context.Background()); rErr != nil {
				logger.FromContext(ctx).Error().Err(rErr).Msg("SecretsRepo - Create - tx.Rollback")
			}
		}
	}()

	err = tx.QueryRow(
		ctx,
		`INSERT
         INTO secrets (owner_id, name, kind, metadata, data)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING secret_id`,
		owner,
		name,
		kind,
		metadata,
		data,
	).Scan(&id)
	if err != nil {
		if postgres.IsEntityExists(err) {
			return id, entity.ErrSecretExists
		}

		return id, fmt.Errorf("SecretsRepo - Create - tx.QueryRow.Scan: %w", err)
	}

	return id, nil
}

// List returns all secrets of the provided user.
func (r *SecretsRepo) List(
	ctx context.Context,
	owner uuid.UUID,
) ([]entity.Secret, error) {
	rv := make([]entity.Secret, 0)
	if err := r.pg.Select(
		ctx,
		&rv,
		`SELECT
         secret_id, name, kind, metadata
     FROM secrets
     WHERE owner_id = $1`,
		owner.String(),
	); err != nil {
		return nil, fmt.Errorf("SecretsRepo - List - r.Select: %w", err)
	}

	return rv, nil
}
