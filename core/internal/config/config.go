package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("empty config path")
	}

	cfg := &Config{}
	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

type Config struct {
	Env        string     `yaml:"env" env-default:"development"`
	DB         DB         `yaml:"db" env-required:"true"`
	HTTPServer HTTPServer `yaml:"httpserver"`
}

type HTTPServer struct {
	Port int `yaml:"int" env-default:"3000"`
	//	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
	//	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	//	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DB struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Name     string `yaml:"name" env-required:"true"`
}
