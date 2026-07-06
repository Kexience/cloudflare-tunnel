package config

import "go.uber.org/fx"

// 注意：configPath 由 main 函数传入
func NewConfigModule(configPath string) fx.Option {
	return fx.Module(
		"config",
		fx.Provide(
			func() (*Config, error) {
				if configPath == "" {
					configPath = "./config"
				}
				return LoadConfig(configPath)
			},
		),
	)
}
