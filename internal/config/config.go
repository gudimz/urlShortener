package config

import (
	"github.com/gudimz/urlShortener/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	Listen struct {
		Ip   string `yaml:"ip" env-default:"127.0.0.1"`
		Port string `yaml:"port" env-default:"8080"`
	} `yaml:"listen"`
}

var (
	instance *Config
	once     sync.Once
)

const pathToConf = "config.yml"

func GetConfig() *Config {
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
