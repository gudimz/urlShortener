package config

import (
	"github.com/gudimz/urlShortener/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

const pathToConfig = "config.yml"

type ServerConfig struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"8080"`
	Protocol string `yaml:"protocol" env-default:"http"`
}

type PostgresConfig struct {
	Username string `yaml:"username" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres"`
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	DbName   string `yaml:"db_name" env-default:"url_shorten"`
}

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Postgres PostgresConfig `yaml:"postgres"`
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Infoln("reading config.yml")

		instance = &Config{}
		err := cleanenv.ReadConfig(pathToConfig, instance)
		if err != nil {
			description, _ := cleanenv.GetDescription(instance, nil)
			logger.Infoln(description)
			logger.Fatalln(err)
		}
	})
	return instance
}
