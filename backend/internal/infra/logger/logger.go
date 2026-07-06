package logger

import (
	"cloudflared-tunnel/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	sugar *zap.SugaredLogger
}

func NewLogger(cfg *config.Config) (Logger, error) {
	var zapCfg zap.Config

	if cfg.App.Env == "production" || cfg.App.Env == "prod" {
		zapCfg = zap.NewProductionConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	} else {
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	zapLogger, err := zapCfg.Build(zap.AddCallerSkip(0))
	if err != nil {
		return Logger{}, err
	}
	return Logger{sugar: zapLogger.Sugar()}, nil
}

func (l Logger) With(args ...any) Logger {
	return Logger{sugar: l.sugar.With(args...)}
}

func (l Logger) Debug(msg string, args ...any) {
	l.sugar.Debugw(msg, args...)
}

func (l Logger) Info(msg string, args ...any) {
	l.sugar.Infow(msg, args...)
}

func (l Logger) Warn(msg string, args ...any) {
	l.sugar.Warnw(msg, args...)
}

func (l Logger) Error(msg string, args ...any) {
	l.sugar.Errorw(msg, args...)
}

func (l Logger) Sync() error {
	return l.sugar.Sync()
}

func NewLoggerForTest() (Logger, error) {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		return Logger{}, err
	}
	return Logger{sugar: zapLogger.Sugar()}, nil
}
