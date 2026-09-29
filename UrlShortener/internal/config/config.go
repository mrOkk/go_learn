package config

import (
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type AppConfig struct {
	GrpcPort  string         `env:"GRPC_PORT"`
	LocalAddr string         `env:"LOCAL_ADDR"`
	LbAddr    string         `env:"LB_ADDR"`
	Redis     RedisConfig    `envPrefix:"REDIS_"`
	Postgres  PostgresConfig `envPrefix:"POSTGRES_"`
}

type RedisConfig struct {
	TTL  time.Duration `env:"TTL" envDefault:"24h"`
	Host string        `env:"HOST"`
	Port string        `env:"PORT"`
}

type PostgresConfig struct {
	Host string `env:"HOST"`
	Port int    `env:"PORT"`
	DB   string `env:"DB"`
	User string `env:"USER"`
	Pass string `env:"PASS"`
	SSL  string `env:"SSL"`
}

func (r RedisConfig) GetOpts() *redis.Options {
	return &redis.Options{
		Addr: fmt.Sprintf("%s:%s", r.Host, r.Port),
	}
}

func (p PostgresConfig) GetConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", p.Host, p.Port, p.User, p.Pass, p.DB, p.SSL)
}

func Load() AppConfig {
	_ = godotenv.Load()
	var config AppConfig
	err := env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("config: %v", config)
	return config
}
