package http

import "go.uber.org/fx"

func NewHttpModule() fx.Option {
	return fx.Module(
		"http",
		fx.Provide(
			NewRouter,
		),
	)
}
