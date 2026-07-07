package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	if configPath == "" {
		configPath = "./config"
	}

	// 设置默认值
	v.SetDefault("app.name", "tunnel")
	v.SetDefault("app.env", "production")
	v.SetDefault("app.port", 8083)
	v.SetDefault("db.driver", "sqlite3")
	v.SetDefault("db.dsn", "file:cloudflared-tunnel.db?cache=shared&_fk=1")
	v.SetDefault("jwt.secret", "")
	v.SetDefault("jwt.expire_hour", 24)
	v.SetDefault("credential.secret", "")
	v.SetDefault("cloudflared.version", "")

	if stat, err := os.Stat(configPath); err == nil {
		if stat.IsDir() {
			v.SetConfigName("config")
			v.SetConfigType("yaml")
			v.AddConfigPath(configPath)
		} else {
			v.SetConfigFile(configPath)
		}
	} else {
		ext := strings.ToLower(filepath.Ext(configPath))
		if ext != "" {
			v.SetConfigFile(configPath)
		} else {
			v.SetConfigName("config")
			v.SetConfigType("yaml")
			v.AddConfigPath(configPath)
		}
	}

	v.AddConfigPath(".")
	v.AutomaticEnv()

	// env key 替换, 例如: APP_NAME -> app.name
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := v.ReadInConfig(); err != nil {
		// 配置文件不存在时仅警告，允许通过环境变量加载配置
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("警告: 未找到配置文件，将使用默认值和环境变量")
		} else {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &cfg, nil
}
