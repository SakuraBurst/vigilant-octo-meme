package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"path/filepath"
	"time"
)

type Config struct {
	Env         string    `yaml:"env" env-default:"local"`
	StoragePath string    `yaml:"storage_path" env-required:"true"`
	App         AppConfig `yaml:"app"`
}

type AppConfig struct {
	Port     string        `yaml:"port"`
	CacheTTL time.Duration `yaml:"cache_ttl"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadConfig(filepath.Join(".", "config", "config.yaml"), cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
