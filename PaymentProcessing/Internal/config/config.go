package config

type Config struct {
	Cache    CacheConfig
	Postgres PostgresConfig
	Kafka    KafkaConfig
}

type AppConfig struct {
	Env string
}

type CacheConfig struct {
	Capacity int
}

type PostgresConfig struct {
	Host string
	Port int
	DB   string
	User string
	Pass string
	SSL  string
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupId string
}

func (a *AppConfig) isDebugMode() bool {
	return a.Env == "debug"
}

func debugCacheConfig() CacheConfig {
	return CacheConfig{Capacity: 1000}
}

func debugPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host: "localhost",
		Port: 5432,
		DB:   "transactions",
		User: "postgres",
		Pass: "password",
		SSL:  "disable",
	}
}

func debugKafkaConfig() KafkaConfig {
	return KafkaConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "transactions",
		GroupId: "payment-consumer-group",
	}
}
