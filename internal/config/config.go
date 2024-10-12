package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

// Config struct for config file
type Config struct {
	Env         string     `yaml:"env" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server"`
	MetricsAddr string     `yaml:"metrics_address"`
	PostgreSQL  PostgreSQL `yaml:"postgresql"`
}

type HTTPServer struct {
	Addr        string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type PostgreSQL struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// MustLoad loads config file
//
// # Panic on error
//
// Returning *Config
func MustLoad() *Config {
	// load .env
	godotenv.Load()

	// Get config path
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	// load config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	return &cfg
}
