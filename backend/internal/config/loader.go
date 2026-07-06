package config

import (
	"fmt"
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
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &cfg, nil
}
