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
	Port     string         `env:"PORT" envDefault:"8080"`
	Redis    RedisConfig    `envPrefix:"REDIS_"`
	Postgres PostgresConfig `envPrefix:"POSTGRES_"`
}

type RedisConfig struct {
	TTL  time.Duration `env:"TTL" envDefault:"24h"`
	host string        `env:"HOST"`
	port string        `env:"PORT"`
}

type PostgresConfig struct {
	host string `env:"HOST"`
	port int    `env:"PORT"`
	db   string `env:"DB"`
	user string `env:"USER"`
	pass string `env:"PASS"`
	ssl  string `env:"SSL"`
}

func (r RedisConfig) GetOpts() *redis.Options {
	return &redis.Options{
		Addr: fmt.Sprintf("%s:%s", r.host, r.port),
	}
}

func (p PostgresConfig) GetConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", p.host, p.port, p.user, p.pass, p.db, p.ssl)
}

func Load() AppConfig {
	_ = godotenv.Load()
	var config AppConfig
	err := env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
