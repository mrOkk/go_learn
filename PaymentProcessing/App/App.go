package main

import (
	"App/Internal/Consumer"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	brokers := []string{"localhost:9092"}
	topic := "payments"
	groupId := "payment-consumer-group"

	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupId, config)
	if err != nil {
		log.Fatal("Error creating consumer group: ", err)
	}
	defer consumerGroup.Close()

	log.Println("Connected to Kafka brokers: ", brokers)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := &Consumer.Handler{}

	go func() {
		select {
		case <-ctx.Done():
			log.Println("Shutting down consumer...")
			return
		default:
		}

		err := consumerGroup.Consume(ctx, []string{topic}, handler)
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
