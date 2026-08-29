package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Address string `yaml:"address"`
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
}
type Config struct {
	Env         string     `yaml:"env"  env:"ENV" envDefault:"production" envRequired:"true"`
	StoragePath string     `yaml:"storage_path" env:"STORAGE_PATH" envDefault:"storage/storage.db" envRequired:"true"`
	HttpServer  HttpServer `yaml:"http_server"`
}

func MustLoadConfig() *Config {
	var cfg Config
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")
	if configPath == "" {
		flags := flag.String("config", "", "Path to the config file")
		flag.Parse()
		configPath = *flags
		if configPath == "" {
			log.Fatal("Config path is not provided. Please set CONFIG_PATH environment variable or use --config flag.")
		}
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist at path: %s", configPath)
	}
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatal("Failed to read config file:", err)
	}
	return &cfg
}
