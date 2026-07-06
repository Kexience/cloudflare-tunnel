package main

import (
	"cloudflared-tunnel/internal/config"
	httpEntry "cloudflared-tunnel/internal/entry/http"
	"cloudflared-tunnel/internal/infra"
	"cloudflared-tunnel/internal/infra/logger"
	"cloudflared-tunnel/internal/module"
	"context"
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func main() {
	var path string
	flag.StringVar(&path, "path", "", "path to serve")
	flag.Parse()

	fx.New(
		config.NewConfigModule(path),
		infra.Module,
		module.Module,
		httpEntry.NewHttpModule(),
		fx.Invoke(StartHttpServer),
	).Run()
}

func StartHttpServer(lc fx.Lifecycle, cfg *config.Config, g *gin.Engine, log logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			host := fmt.Sprintf("%s:%v", "0.0.0.0", cfg.App.Port)
			log.Info("开始监听", "host", host)
			go g.Run(host)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
