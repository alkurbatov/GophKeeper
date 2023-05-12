package repo

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/alkurbatov/goph-keeper/internal/keeper/entity"
	"github.com/alkurbatov/goph-keeper/internal/keeper/infra/postgres"
	"github.com/alkurbatov/goph-keeper/pkg/goph"
	uuid "github.com/satori/go.uuid"
)

var _ Secrets = (*SecretsRepo)(nil)

var ErrNoValuesToUpdate = errors.New("no values to update")

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
) (id uuid.UUID, err error) {
	fn := func(tx postgres.Transaction) error {
		err := tx.QueryRow(
			ctx,
			`INSERT INTO
           secrets (owner_id, name, kind, metadata, data)
       VALUES
           ($1, $2, $3, $4, $5)
       RETURNING secret_id`,
			owner,
			name,
			kind,
			metadata,
			data,
		).Scan(&id)
		if err != nil {
			if postgres.IsEntityExists(err) {
				return entity.ErrSecretExists
			}

			return fmt.Errorf("SecretsRepo - Create - tx.QueryRow.Scan: %w", err)
		}

		return nil
	}

	if err := r.pg.RunAtomic(ctx, fn); err != nil {
		return id, fmt.Errorf("SecretsRepo - Create - r.pg.RunAtomic: %w", err)
	}

	return id, nil
}

// List returns all secrets of the provided user.
// Data is not filled in this case to reduce load on service.
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
     FROM
         secrets
     WHERE owner_id = $1`,
		owner,
	); err != nil {
		return nil, fmt.Errorf("SecretsRepo - List - r.Select: %w", err)
	}

	return rv, nil
}

// Get returns full secret info and data.
func (r *SecretsRepo) Get(
	ctx context.Context,
	owner, id uuid.UUID,
) (*entity.Secret, error) {
	var secret entity.Secret

	err := r.pg.Pool.
		QueryRow(
			ctx,
			`SELECT
           secret_id, name, kind, metadata, data
       FROM
           secrets
       WHERE secret_id=$1 AND owner_id = $2`,
			id,
			owner,
		).
		Scan(&secret.ID, &secret.Name, &secret.Kind, &secret.Metadata, &secret.Data)
	if err != nil {
		if postgres.IsEmptyResponse(err) {
			return nil, entity.ErrSecretNotFound
		}

		return nil, fmt.Errorf("SecretsRepo - Get - r.pg.Pool.QueryRow.Scan: %w", err)
	}

	return &secret, nil
}

// Update changes secret info and data.
func (r *SecretsRepo) Update( //nolint:cyclop,gocognit,gocyclo // update is complex operation
	ctx context.Context,
	owner, id uuid.UUID,
	changed []string,
	name string,
	metadata, data []byte,
) error {
	fn := func(tx postgres.Transaction) error {
		values := []any{}
		query := "UPDATE secrets SET"

		for _, field := range changed {
			if len(values) != 0 {
				query += ","
			}

			switch field {
			case "name":
				values = append(values, name)
				query += " name = $" + strconv.Itoa(len(values))

			case "metadata":
				values = append(values, metadata)
				query += " metadata = $" + strconv.Itoa(len(values))

			case "data":
				values = append(values, data)
				query += " data = $" + strconv.Itoa(len(values))
			}
		}

		if len(values) == 0 {
			return fmt.Errorf("SecretsRepo - Update: %w", ErrNoValuesToUpdate)
		}

		values = append(values, id)
		query += " WHERE secret_id = $" + strconv.Itoa(len(values))

		values = append(values, owner)
		query += " AND owner_id = $" + strconv.Itoa(len(values))

		tag, err := tx.Exec(
			ctx,
			query,
			values...,
		)
		if err != nil {
			if postgres.IsEntityExists(err) {
				return entity.ErrSecretNameConflict
			}

			return fmt.Errorf("SecretsRepo - Update - tx.Exec: %w", err)
		}

		if tag.RowsAffected() == 0 {
			return entity.ErrSecretNotFound
		}

		return nil
	}

	if err := r.pg.RunAtomic(ctx, fn); err != nil {
		return fmt.Errorf("SecretsRepo - Update - r.pg.RunAtomic: %w", err)
	}

	return nil
}

// Delete removes secret from database.
func (r *SecretsRepo) Delete(
	ctx context.Context,
	owner, id uuid.UUID,
) (err error) {
	fn := func(tx postgres.Transaction) error {
		tag, err := tx.Exec(
			ctx,
			`DELETE FROM
           secrets
       WHERE secret_id = $1 AND owner_id = $2`,
			id,
			owner,
		)
		if err != nil {
			return fmt.Errorf("SecretsRepo - Delete - tx.Exec: %w", err)
		}

		if tag.RowsAffected() == 0 {
			return entity.ErrSecretNotFound
		}

		return nil
	}

	if err := r.pg.RunAtomic(ctx, fn); err != nil {
		return fmt.Errorf("SecretsRepo - Delete - r.pg.RunAtomic: %w", err)
	}

	return nil
}
