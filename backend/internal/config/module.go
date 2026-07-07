package config

import (
	"fmt"

	"cloudflared-tunnel/pkg/crypto"

	"go.uber.org/fx"
)

// 注意：configPath 由 main 函数传入
func NewConfigModule(configPath string) fx.Option {
	return fx.Module(
		"config",
		fx.Provide(
			func() (*Config, error) {
				if configPath == "" {
					configPath = "./config"
				}
				cfg, err := LoadConfig(configPath)
				if err != nil {
					return nil, err
				}
				if err := validate(cfg); err != nil {
					return nil, err
				}
				return cfg, nil
			},
		),
	)
}

func validate(cfg *Config) error {
	if len(cfg.Credential.Secret) != crypto.KeyLength {
		return fmt.Errorf("credential.secret 长度必须为 %d 字节，当前为 %d 字节", crypto.KeyLength, len(cfg.Credential.Secret))
	}
	if cfg.JWT.Secret == "" {
		return fmt.Errorf("jwt.secret 未设置，请通过环境变量 JWT_SECRET 或配置文件设置")
	}
	return nil
}
