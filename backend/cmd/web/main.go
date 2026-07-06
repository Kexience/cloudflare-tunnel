package main

import (
	"cloudflared-tunnel/internal/config"
	httpEntry "cloudflared-tunnel/internal/entry/http"
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
		module.Module,
		httpEntry.NewHttpModule(),
		fx.Invoke(StartHttpServer),
	).Run()
}

func StartHttpServer(lc fx.Lifecycle, cfg *config.Config, g *gin.Engine) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			host := fmt.Sprintf("%s:%v", "0.0.0.0", cfg.App.Port)
			go g.Run(host)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
