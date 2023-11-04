package ds

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	Host    string `yaml:"host" env-default:"0.0.0.0"`
	Port    string `yaml:"port" env-default:"8080"`
	BaseUrl string `yaml:"base_url" env-default:"http://localhost"`
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

const (
	configPathEnv = "CONFIG_PATH"
)

func GetConfig() *Config {
	configPath := os.Getenv(configPathEnv)
	if configPath == "" {
		log.Fatal("env CONFIG_PATH not set")
	}

	log.Printf("reading config file: %s", configPath)
	config := Config{}
	err := cleanenv.ReadConfig(configPath, &config)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	return &config
}
