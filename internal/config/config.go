package config

import (
	"github.com/gudimz/urlShortener/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type ListenConfig struct {
	Ip   string `yaml:"ip" env-default:"127.0.0.1"`
	Port string `yaml:"port" env-default:"8080"`
}

type PostgresConfig struct {
	Username string `yaml:"username" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres"`
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	DbName   string `yaml:"db_name" env-default:"url_shorten"`
}

type Config struct {
	Listen   ListenConfig   `yaml:"listen"`
	Postgres PostgresConfig `yaml:"postgres"`
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig(pathToConf string) *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Infoln("reading config.yml")

		instance = &Config{}
		err := cleanenv.ReadConfig(pathToConf, instance)
		if err != nil {
			description, _ := cleanenv.GetDescription(instance, nil)
			logger.Infoln(description)
			logger.Fatalln(err)
		}
	})
	return instance
}
