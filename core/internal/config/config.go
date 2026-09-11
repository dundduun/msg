package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
	"time"
)

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("empty config path")
	}

	envPath := os.Getenv("DOTENV_PATH")
	if envPath != "" {
		err := godotenv.Load(envPath)
		if err != nil {
			panic("failed to load .env: " + err.Error())
		}
	}

	cfg := &Config{}
	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

type Config struct {
	Env        string     `env:"ENV" yaml:"env" env-default:"development"`
	DB         DB         `yaml:"db"`
	Cache      Cache      `yaml:"cache"`
	HTTPServer HTTPServer `yaml:"httpserver"`
}

type DB struct {
	Host     string `env:"DB_HOST" yaml:"host" env-required:"true"`
	Port     int    `env:"DB_PORT" yaml:"port" env-required:"true"`
	User     string `env:"DB_USER" yaml:"user" env-required:"true"`
	Password string `env:"DB_PASSWORD" yaml:"password" env-required:"true"`
	Name     string `env:"DB_NAME" yaml:"name" env-required:"true"`
}

type Cache struct {
	Host     string        `env:"REDIS_HOST" yaml:"host" env-required:"true"`
	Port     int           `env:"REDIS_PORT" yaml:"port" env-required:"true"`
	Password string        `env:"REDIS_PASSWORD" yaml:"password" env-required:"true"`
	Name     int           `env:"REDIS_NAME" yaml:"name" env-required:"true"`
	TTL      time.Duration `yaml:"ttl" env-required:"true"`
}

type HTTPServer struct {
	Port int `yaml:"port"`
	//	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
	//	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	//	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}
