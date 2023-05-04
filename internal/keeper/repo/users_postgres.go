package repo

import (
	"context"
	"fmt"

	"github.com/alkurbatov/goph-keeper/internal/keeper/entity"
	"github.com/alkurbatov/goph-keeper/internal/keeper/infra/logger"
	"github.com/alkurbatov/goph-keeper/internal/keeper/infra/postgres"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

var _ Users = (*UsersRepo)(nil)

// UsersRepo is facade to users stored in Postgres.
type UsersRepo struct {
	pg  *postgres.Postgres
	log *logger.Logger
}

// NewUsersRepo creates and initializes UsersRepo object.
func NewUsersRepo(
	pg *postgres.Postgres,
	log *logger.Logger,
) *UsersRepo {
	return &UsersRepo{pg, log}
}

// Register creates a new user.
func (r *UsersRepo) Register(
	ctx context.Context,
	username, securityKey string,
) (id uuid.UUID, err error) {
	tx, err := r.pg.BeginTx(ctx)
	if err != nil {
		return id, fmt.Errorf("UsersRepo - Register - r.pg.BeginTx: %w", err)
	}

	defer func() {
		switch err {
		case nil:
			if cErr := tx.Commit(context.Background()); cErr != nil {
				err = fmt.Errorf("UsersRepo - Register - tx.Commit: %w", cErr)
			}
		default:
			if rErr := tx.Rollback(context.Background()); rErr != nil {
				log.Ctx(ctx).Error().Err(rErr).Msg("UsersRepo - Register - tx.Rollback")
			}
		}
	}()

	err = tx.QueryRow(
		ctx,
		`INSERT
         INTO users (username, security_key)
         VALUES ($1, crypt($2, gen_salt('bf', 8)))
         RETURNING user_id`,
		username,
		securityKey,
	).Scan(&id)
	if err != nil {
		if postgres.IsEntityExists(err) {
			return id, entity.ErrUserExists
		}

		return id, fmt.Errorf("UsersRepo - Register - tx.QueryRow.Scan: %w", err)
	}

	return id, nil
}

// Verify checks provided username and security key against data stored in database.
// Returns entity.User, if verification was successful.
func (r *UsersRepo) Verify(
	ctx context.Context,
	username, securityKey string,
) (entity.User, error) {
	var user entity.User

	err := r.pg.Pool.
		QueryRow(
			ctx,
			`SELECT user_id, username FROM users
           WHERE username=$1 AND security_key = crypt($2, security_key)`,
			username,
			securityKey,
		).
		Scan(&user.ID, &user.Username)
	if err != nil {
		if postgres.IsEmptyResponse(err) {
			return user, entity.ErrInvalidCredentials
		}

		return user, fmt.Errorf("UsersRepo - Verify - r.pg.Pool.QueryRow.Scan: %w", err)
	}

	return user, nil
}
