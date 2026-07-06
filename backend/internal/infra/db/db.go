package db

import (
	"cloudflared-tunnel/ent"
	"cloudflared-tunnel/internal/config"
	"context"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewClient(lc fx.Lifecycle, cfg *config.Config, logger *zap.Logger) (*ent.Client, error) {
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
			logger.Info("Running database migrations")
			if err := client.Schema.Create(ctx); err != nil {
				return fmt.Errorf("failed creating schema resources: %w", err)
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing database connection")
			return client.Close()
		},
	})

	return client, nil
}
