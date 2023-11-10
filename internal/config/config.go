package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-required:"true"`
	StoragePath string `yaml:"storage_path"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`

	User     string `yaml:"user" env-required:"true"`
	Password string `password:"user" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

func MustLoad(configPath string) Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}
	var cnfg Config
	if err := cleanenv.ReadConfig(configPath, &cnfg); err != nil {
		log.Fatalf("cannot read config: %s", configPath)
	}

	return cnfg
}
