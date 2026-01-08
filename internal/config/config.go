package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env:"ENV" env-default:"local" env-required:"true"`
	StorageUrl string `yaml:"storage_url" env-required:"true"`
	HttpServer `yaml:"http_server"`
	Security   `yaml:"security"`
}

type HttpServer struct {
	Address     string        `yaml: "address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml: "timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml: "idle_timeout" env-default:"10s"`
}

type Security struct {
	TokenTTL   int64 `yaml:"token_ttl" env-default:"3600"`
	SigningKey string
	HasherSalt string
}

func MustLoad() Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}

	setSigningKey(cfg)

	return cfg
}

func setSigningKey(cfg Config) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	signingKey := os.Getenv("SIGNING_KEY")
	if signingKey == "" {
		log.Fatal("SIGNING_KEY environment variable not set")
	}

	hasherSalt := os.Getenv("SIGNING_KEY_SALT")
	if hasherSalt == "" {
		log.Fatal("SIGNING_KEY_SALT environment variable not set")
	}

	cfg.Security.SigningKey = signingKey
	cfg.Security.HasherSalt = hasherSalt
}
