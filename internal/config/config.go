package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	DB         `yaml:"db"`
	HttpServer `yaml:"http_server"`
}

type HttpServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DB struct {
	Host     string `yaml:"host" env-default:"localhost" env:"DB_HOST"`
	Port     string `yaml:"port" env-default:"5432" env:"DB_PORT"`
	User     string `yaml:"user" env-required:"true" env:"DB_USER"`
	Password string `yaml:"password" env-required:"true" env:"DB_PASSWORD"`
	DBName   string `yaml:"dbName" env-required:"true" env:"DB_NAME"`
}

func MustLoad() *Config {
	cfgPath := os.Getenv("CONFIG_PATH")

	if len(strings.TrimSpace(cfgPath)) == 0 {
		log.Fatal("CONFIG_PATH not set")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatal("file not found")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		log.Fatalf("failed read config:%s", err.Error())
	}

	return &cfg
}
