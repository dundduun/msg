package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("empty config path")
	}
	var cfg *Config
	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

type Config struct {
	Env string `yaml:"env" env-default:"development"`
	HTTPServer
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}
