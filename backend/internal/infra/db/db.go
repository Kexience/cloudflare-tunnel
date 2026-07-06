package db

import (
	"cloudflared-tunnel/ent"
	"cloudflared-tunnel/internal/config"
	"cloudflared-tunnel/internal/infra/logger"
	"context"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/fx"
)

func NewClient(lc fx.Lifecycle, cfg *config.Config, log logger.Logger) (*ent.Client, error) {
	dbCfg := cfg.DB
	if dbCfg.Driver == "" {
		dbCfg.Driver = "sqlite3"
	}
	if dbCfg.DSN == "" {
		dbCfg.DSN = "file:cloudflared-tunnel.db?cache=shared&_fk=1"
	}

	client, err := ent.Open(dbCfg.Driver, dbCfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to sqlite: %w", err)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Running database migrations")
			if err := client.Schema.Create(ctx); err != nil {
				return fmt.Errorf("failed creating schema resources: %w", err)
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Closing database connection")
			return client.Close()
		},
	})

	return client, nil
}
