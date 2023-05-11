// Package app implements keeper service.
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alkurbatov/goph-keeper/internal/keeper/config"
	v1 "github.com/alkurbatov/goph-keeper/internal/keeper/controller/grpc/v1"
	"github.com/alkurbatov/goph-keeper/internal/keeper/infra/grpcserver"
	"github.com/alkurbatov/goph-keeper/internal/keeper/infra/logger"
	"github.com/alkurbatov/goph-keeper/internal/keeper/infra/postgres"
	"github.com/alkurbatov/goph-keeper/internal/keeper/repo"
	"github.com/alkurbatov/goph-keeper/internal/keeper/usecase"
	"google.golang.org/grpc"
)

const (
	_defaultMinimalSecretLength = 32
	_defaultShutdownTimeout     = 60 * time.Second
)

// Run initializes and starts the keeper service.
func Run(cfg *config.Config) error { //nolint:funlen // no sense to split main Run function
	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("app - Run - logger.New: %w", err)
	}

	log.Info().Msg(cfg.String())

	if len([]byte(cfg.Secret)) < _defaultMinimalSecretLength {
		log.Warn().Msg("Insecure signature: secret key is shorter than 32 bytes!")
	}

	pg, err := postgres.New(string(cfg.DatabaseURI), log)
	if err != nil {
		return fmt.Errorf("app - Run - postgres.New: %w", err)
	}

	repos := repo.New(pg)
	usecases := usecase.New(cfg, repos)

	grpcSrv, err := grpcserver.New(
		cfg.Address,
		cfg.CrtPath,
		cfg.KeyPath,
		grpc.MaxRecvMsgSize(v1.DefaultMaxMessageSize),
		grpc.ChainUnaryInterceptor(
			v1.LoggingUnarysInterceptor(log),
			v1.AuthUnaryInterceptor(cfg.Secret),
		),
	)
	if err != nil {
		return fmt.Errorf("app - Run - grpcserver.New: %w", err)
	}

	v1.RegisterRoutes(grpcSrv.Instance(), usecases)
	grpcSrv.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	select {
	case s := <-interrupt:
		log.Info().Msg("app - Run - interrupt: signal " + s.String())

	case err := <-grpcSrv.Notify():
		log.Error().Err(err).Msg("app - Run - grpcSrv.Notify")
	}

	log.Info().Msg("Shutting down...")

	stopped := make(chan struct{})

	stopCtx, cancel := context.WithTimeout(context.Background(), _defaultShutdownTimeout)
	defer cancel()

	go func() {
		shutdown(log, grpcSrv, pg)
		close(stopped)
	}()

	select {
	case <-stopped:
		log.Info().Msg("Service shutdown successful")

	case <-stopCtx.Done():
		log.Warn().Msgf("Exceeded %s shutdown timeout, exit forcibly", _defaultShutdownTimeout)
	}

	return nil
}

func shutdown(
	log *logger.Logger,
	grpcSrv *grpcserver.Server,
	pg *postgres.Postgres,
) {
	log.Info().Msg("Shutting down gRPC API...")
	grpcSrv.Shutdown()

	log.Info().Msg("Shutting down database connection...")
	pg.Close()
}
