package config

type Config struct {
	App        AppConfig        `mapstructure:"app"`
	DB         DBConfig         `mapstructure:"db"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Credential CredentialConfig `mapstructure:"credential"`
	Cloudflared CloudflaredConfig `mapstructure:"cloudflared"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
	Port int    `mapstructure:"port"`
}

type DBConfig struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpireHour int    `mapstructure:"expire_hour"`
}

type CredentialConfig struct {
	Secret string `mapstructure:"secret"`
}

type CloudflaredConfig struct {
	Version string `mapstructure:"version"`
}
