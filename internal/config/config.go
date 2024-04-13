package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"path/filepath"
	"time"
)

type Config struct {
	Env         string    `yaml:"env" env-default:"local" env:"ENVIRONMENT"`
	StoragePath string    `yaml:"storage_path" env-required:"true" env:"STORAGE_PATH"`
	App         AppConfig `yaml:"app"`
}

type AppConfig struct {
	Port      string        `yaml:"port" env:"PORT" env-default:"8080"`
	CacheTTL  time.Duration `yaml:"cache_ttl" env:"CACHE_TTL" env-default:"5m"`
	TokenTTL  time.Duration `yaml:"token_ttl" env:"TOKEN_TTL" env-default:"1h"`
	AppSecret string        `yaml:"app_secret" env:"APP_SECRET" env-required:"true"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadConfig(filepath.Join(".", "config", "config.yaml"), cfg)
	if err != nil {
		err = cleanenv.ReadEnv(cfg)
		if err != nil {
			panic(err)
		}
	}
	return cfg
}
