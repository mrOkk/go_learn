package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Load() (Config, error) {
	app, _ := loadAppConfig()
	if app.isDebugMode() {
		return Config{
			Cache:    debugCacheConfig(),
			Postgres: debugPostgresConfig(),
			Kafka:    debugKafkaConfig(),
			Worker:   debugWorkerConfig(),
		}, nil
	}

	cache, err := loadCacheConfig()
	if err != nil {
		return Config{}, err
	}

	postgres, err := loadPostgresConfig()
	if err != nil {
		return Config{}, err
	}

	kafka, err := loadKafkaConfig()
	if err != nil {
		return Config{}, err
	}

	worker, err := loadWorkerConfig()
	if err != nil {
		return Config{}, err
	}

	return Config{
		Cache:    cache,
		Postgres: postgres,
		Kafka:    kafka,
		Worker:   worker,
	}, nil
}

func loadWorkerConfig() (WorkerConfig, error) {
	// TODO
	return debugWorkerConfig(), nil
}

func loadAppConfig() (AppConfig, error) {
	env, err := getString("APP_ENV")
	if err != nil {
		return AppConfig{}, err
	}

	return AppConfig{
		Env: env,
	}, nil
}

func loadCacheConfig() (CacheConfig, error) {
	capacity, err := getInt("CACHE_CAPACITY")
	if err != nil {
		return CacheConfig{}, err
	}
	if capacity <= 0 {
		return CacheConfig{}, fmt.Errorf("CACHE_CAPACITY must be greater than 0, got %d", capacity)
	}

	return CacheConfig{
		Capacity: capacity,
	}, nil
}

func loadPostgresConfig() (PostgresConfig, error) {
	host, err := getString("POSTGRES_HOST")
	if err != nil {
		return PostgresConfig{}, err
	}
	port, err := getInt("POSTGRES_PORT")
	if err != nil {
		return PostgresConfig{}, err
	}
	if port <= 0 {
		return PostgresConfig{}, fmt.Errorf("POSTGRES_PORT must be greater than 0, got %d", port)
	}
	db, err := getString("POSTGRES_DB")
	if err != nil {
		return PostgresConfig{}, err
	}
	user, err := getString("POSTGRES_USER")
	if err != nil {
		return PostgresConfig{}, err
	}
	password, err := getString("POSTGRES_PASSWORD")
	if err != nil {
		return PostgresConfig{}, err
	}
	sslMode, err := getString("POSTGRES_SSLMODE")
	if err != nil {
		return PostgresConfig{}, err
	}

	return PostgresConfig{
		Host: host,
		Port: port,
		DB:   db,
		User: user,
		Pass: password,
		SSL:  sslMode,
	}, nil
}

func loadKafkaConfig() (KafkaConfig, error) {
	brokers, err := getCSV("KAFKA_BROKERS")
	if err != nil {
		return KafkaConfig{}, err
	}
	topic, err := getString("KAFKA_TOPIC")
	if err != nil {
		return KafkaConfig{}, err
	}
	groupID, err := getString("KAFKA_GROUP_ID")
	if err != nil {
		return KafkaConfig{}, err
	}

	return KafkaConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupId: groupID,
	}, nil
}

func getString(key string) (string, error) {
	raw, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("environment variable %s not set", key)
	}
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", fmt.Errorf("%s must not be empty", key)
	}
	return value, nil
}

func getInt(key string) (int, error) {
	raw, err := getString(key)
	if err != nil {
		return 0, err
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valud integer, got %q: %w", key, raw, err)
	}
	return n, nil
}

func getCSV(key string) ([]string, error) {
	raw, err := getString(key)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(raw, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item == "" {
			return nil, fmt.Errorf("%s contains empty item in %q", key, raw)
		}
		values = append(values, item)
	}

	if len(values) == 0 {
		return nil, fmt.Errorf("%s must contain at least one value", key)
	}

	return values, nil
}
