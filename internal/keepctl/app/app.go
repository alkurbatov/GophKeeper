// Package app implements keepctl service.
package app

import (
	"context"
	"errors"

	"github.com/alkurbatov/goph-keeper/internal/keepctl/config"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/entity"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/infra/grpcconn"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/infra/logger"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/repo"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/usecase"
)

type appKey string

var (
	appKeyName appKey = "app"

	errNotInitialized = errors.New("application is not initialized")
)

// App implements client application for the keeper service.
type App struct {
	Log         *logger.Logger
	conn        *grpcconn.Connection
	Usecases    *usecase.UseCases
	Key         entity.Key
	AccessToken string
}

// New creates and initializes new App object.
func New(cfg *config.Config) (*App, error) {
	log := logger.New(cfg.Verbose)
	log.Debug().Msg(cfg.String())

	conn, err := grpcconn.New(cfg.Address, cfg.CAPath)
	if err != nil {
		log.Debug().Err(err).Msg("app - New - grpcconn.New")

		return nil, err
	}

	key := entity.NewKey(cfg.Username, cfg.Password)
	repos := repo.New(conn)
	usecases := usecase.New(key, repos)

	return &App{
		Log:      log,
		conn:     conn,
		Usecases: usecases,
	}, nil
}

// WithContext injects App into provided context.
func (a *App) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, appKeyName, a)
}

// FromContext extracts App from provided context.
func FromContext(ctx context.Context) (*App, error) {
	if val := ctx.Value(appKeyName); val != nil {
		return val.(*App), nil
	}

	return nil, errNotInitialized
}

// Shutdown gracefully stops client application.
func (a *App) Shutdown() {
	if err := a.conn.Close(); err != nil {
		a.Log.Warn().Err(err).Msg("app - Shutdown - a.conn.Close")
	}
}
