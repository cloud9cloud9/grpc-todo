package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"log"
	"sync"
)

type Config struct {
	Port       string `mapstructure:"PORT"`
	AuthSuvURL string `mapstructure:"AUTH_SUV_URL"`
	TodoSuvURL string `mapstructure:"TODO_SUV_URL"`
}

const (
	cfgFile = ".env"
)

var (
	Instance *Config
	once     sync.Once
	envs     = []string{
		"PORT", "AUTH_SUV_URL", "TODO_SUV_URL",
	}
)

func LoadConfig() (*Config, error) {
	once.Do(func() {
		viper.SetConfigFile(cfgFile)

		Instance = &Config{}

		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("Error reading config file, %s", err)
		}

		for _, env := range envs {
			if err := viper.BindEnv(env); err != nil {
				log.Fatalf("Error reading config file, %s", err)
			}
		}

		if err := viper.Unmarshal(Instance); err != nil {
			log.Fatalf("Error unmarshalling config file, %s", err)
		}

		if err := validator.New().Struct(Instance); err != nil {
			log.Fatalf("Error validating config file, %s", err)
		}
	})

	log.Println("Config loaded successfully")
	return Instance, nil
}
