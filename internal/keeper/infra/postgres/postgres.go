package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/alkurbatov/goph-keeper/internal/keeper/infra/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

var DefaultTxOptions = pgx.TxOptions{
	IsoLevel: pgx.ReadCommitted,
}

type PgxIface interface {
	Acquire(ctx context.Context) (*pgxpool.Conn, error)

	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row

	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Close()
}

// Postgres represents abstraction over database connection pool.
type Postgres struct {
	Pool PgxIface
}

// New initializes connection to database and creates new wrapper object.
func New(url string, log *logger.Logger) (*Postgres, error) {
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - New - pgxpool.ParseConfig: %w", err)
	}

	pg := new(Postgres)
	attempts := _defaultConnAttempts

	for attempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Info().Msgf("Postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultConnTimeout)

		attempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - New - attempts == 0: %w", err)
	}

	return pg, nil
}

// Close gracefully closes connection to database.
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

// BeginTx starts new database transaction.
func (p *Postgres) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := p.Pool.BeginTx(ctx, DefaultTxOptions)
	if err != nil {
		return nil, fmt.Errorf("postgres - BeginTx - conn.BeginTx: %w", err)
	}

	return tx, nil
}
