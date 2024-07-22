package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DB DBConfig `yaml:"redis"`
	Bot BotConfig `yaml:"bot"`
	Server ServerConfig `yaml:"server"`
}

type DBConfig struct {
	Host string `yaml:"host" env:"DB_HOST" env-default:"127.0.0.1"`
	Port int `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Username string `yaml:"username" env:"DB_USERNAME" env-default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	DBName string `yaml:"dbname" env:"DB_NAME"`
}

type BotConfig struct {
	Token string `yaml:"token" env:"BOT_TOKEN"`
	UpdateTimeout int `yaml:"timeout" env-default:"10"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int `yaml:"port"`
}

func MustLoad() *Config {
	path := loadPath()
	if path == "" {
		panic("Can`t read config file")
	}

	return loadConfig(path)
}

func loadPath() string {
	var path string
	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()
	if path == "" {	
		path = "./config/config.yaml"
	}

	return path
}

func loadConfig(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}
	
	return &cfg
}