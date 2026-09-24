package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port int        `env:"PORT" envDefault:"8882"`
	GRPC GRPCConfig `envPrefix:"GRPC_"`
}

type GRPCConfig struct {
	Port    int           `env:"PORT" envDefault:"8883"`
	Timeout time.Duration `env:"TIMEOUT" envDefault:"1m"`
}

func LoadConfig() (Config, error) {
	_ = godotenv.Load()
	var config Config
	err := env.Parse(&config)
	return config, err
}
