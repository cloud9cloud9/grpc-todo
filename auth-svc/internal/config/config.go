package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"log"
	"sync"
)

type Config struct {
	Port             string `mapstructure:"PORT"`
	UrlDB            string `mapstructure:"URL_DB"`
	PostgresUser     string `mapstructure:"POSTGRES_AUTH_USER"`
	PostgresPassword string `mapstructure:"POSTGRES_AUTH_PASSWORD"`
	PostgresDB       string `mapstructure:"POSTGRES_AUTH_DB"`
	JwtSecret        string `mapstructure:"JWT_SECRET"`
	JwtIssuer        string `mapstructure:"JWT_ISSUER"`
}

const (
	cfgFile = ".env"
)

var (
	Instance *Config
	once     sync.Once
	envs     = []string{
		"PORT", "URL_DB", "POSTGRES_AUTH_USER",
		"POSTGRES_AUTH_PASSWORD", "POSTGRES_AUTH_DB",
		"JWT_SECRET", "JWT_ISSUER",
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
