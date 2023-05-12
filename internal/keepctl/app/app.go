// Package app implements keepctl service.
package app

import (
	"github.com/alkurbatov/goph-keeper/internal/keepctl/config"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/infra/grpcconn"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/infra/logger"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/repo"
	"github.com/alkurbatov/goph-keeper/internal/keepctl/usecase"
)

// App implements client application for the keeper service.
type App struct {
	Log      *logger.Logger
	conn     *grpcconn.Connection
	Usecases *usecase.UseCases
}

// New creates and initializes new App object.
func New(cfg *config.Config) (*App, error) {
	log := logger.New(cfg.Verbose)

	conn, err := grpcconn.New(cfg.Address, cfg.CAPath)
	if err != nil {
		log.Debug().Err(err).Msg("app - New - grpcconn.New")

		return nil, err
	}

	repos := repo.New(conn)
	usecases := usecase.New(repos)

	return &App{
		Log:      log,
		conn:     conn,
		Usecases: usecases,
	}, nil
}

// Shutdown gracefully stops client application.
func (a *App) Shutdown() {
	if err := a.conn.Close(); err != nil {
		a.Log.Warn().Err(err).Msg("app - Shutdown - a.conn.Close")
	}
}
