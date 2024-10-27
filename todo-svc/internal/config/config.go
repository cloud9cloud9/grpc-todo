package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
)

const (
	cfgPath = "config/config.yml"
)

type Config struct {
	Server struct {
		Port string `yaml:"port" env:"PORT"`
	} `yaml:"server"`

	Database struct {
		Port             string `yaml:"port"`
		Host             string `yaml:"host"`
		PostgresUser     string `yaml:"user" env:"POSTGRES_TODO_USER"`
		PostgresPassword string `yaml:"password" env:"POSTGRES_TODO_PASSWORD"`
		PostgresDB       string `yaml:"database" env:"POSTGRES_TODO_DB"`
		PostgresSSLMode  string `yaml:"sslmode"`
	} `yaml:"db"`
}

var Instance *Config

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Instance = &Config{}

	if err := cleanenv.ReadConfig(cfgPath, Instance); err != nil {
		help, _ := cleanenv.GetDescription(Instance, nil)
		log.Fatalf("Config error: %s", help)
	}

	return Instance, nil
}
