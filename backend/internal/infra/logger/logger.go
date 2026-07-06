package logger

import (
	"cloudflared-tunnel/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	var zapCfg zap.Config

	if cfg.App.Env == "production" || cfg.App.Env == "prod" {
		zapCfg = zap.NewProductionConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	} else {
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	return zapCfg.Build(zap.AddCallerSkip(0))
}
