package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/alkurbatov/goph-keeper/internal/keeper/infra/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Postgres represents abstraction over database connection pool.
type Postgres struct {
	Pool *pgxpool.Pool
	log  *logger.Logger
}

// New initializes connection to database and creates new wrapper object.
func New(url string, log *logger.Logger) (*Postgres, error) {
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - New - pgxpool.ParseConfig: %w", err)
	}

	pg := &Postgres{log: log}
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
