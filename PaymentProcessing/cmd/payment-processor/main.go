package main

import (
	"App/internal/adapters/kafka"
	"App/internal/adapters/postgres"
	"App/internal/application"
	"App/internal/cache"
	"App/internal/config"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	kafkaCfg := cfg.Kafka

	saramaCfg := sarama.NewConfig()
	saramaCfg.Version = sarama.V4_0_0_0
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	saramaCfg.Consumer.Return.Errors = true

	consumerGroup, err := sarama.NewConsumerGroup(kafkaCfg.Brokers, kafkaCfg.GroupId, saramaCfg)
	if err != nil {
		log.Fatal("Error creating consumer group: ", err)
	}
	defer consumerGroup.Close()

	log.Println("Connected to Kafka brokers: ", kafkaCfg.Brokers)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	merchantCache := cache.NewLRUCache(cfg.Cache.Capacity)

	pool, err := pgxpool.New(ctx, cfg.Postgres.BuildConnectionString())
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	defer pool.Close()

	merchantRepo := postgres.NewMerchantRepository(pool)
	txRepo := postgres.NewTransactionRepository(pool)
	txManager := postgres.NewTransactionManager(pool)

	processor := application.NewProcessor(merchantRepo, merchantRepo, txRepo, txManager, merchantCache)
	handler := kafka.NewConsumerHandler(processor, cfg.Worker.QueueSize)

	go func() {
		select {
		case <-ctx.Done():
			log.Println("Shutting down consumer...")
			return
		default:
		}

		err := consumerGroup.Consume(ctx, []string{kafkaCfg.Topic}, handler)
		if err != nil {
			log.Println("Error from consumer: ", err.Error())
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for err := range consumerGroup.Errors() {
			log.Printf("Error from consumer: %v", err.Error())
		}
	}()

	log.Println("Consumer is running. Press Ctrl+C to exit.")

	sinChan := make(chan os.Signal, 1)
	signal.Notify(sinChan, syscall.SIGINT, syscall.SIGTERM)

	<-sinChan

	log.Println("Shutting down consumer...")

	cancel()

	time.Sleep(2 * time.Second)

	log.Println("Consumer shutdown complete.")
}
