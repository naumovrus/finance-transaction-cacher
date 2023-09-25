package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string `yaml:"env" env-default:"local"`
	ConfigDatabase `yaml:"database_config"`
	HTTPServer     `yaml:"http_server"`
}

type HTTPServer struct {
	Port string `yaml:"port" env-default:"8000"`
}

type ConfigDatabase struct {
	Username   string `yaml:"username" env-default:"postgres"`
	PasswordDb string `yaml:"db_password"`
	Host       string `yaml:"host"`
	PortDb     string `yaml:"port_db" env-default:"5436"`
	Dbname     string `yaml:"dbname" env-default:"postgres"`
	Sslmode    string `yaml:"sslmode" env-default:"disable"`
}

func LoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalf("config path is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error occured while loaded conf files: %s", err)
	}
	return &cfg
}
