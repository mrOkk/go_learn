package kafka

import (
	. "App/internal/domain"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type ConsumerHandler struct{}

func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Consumer group session setup")
	return nil
}

func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("Consumer group session cleanup")
	return nil
}

func (h *ConsumerHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Println("Message channel closed")
				return nil
			}

			var transaction Transaction
			if err := json.Unmarshal(message.Value, &transaction); err != nil {
				log.Printf("Failed to decode message at offset %d: %v", message.Offset, err)
				continue
			}

			log.Println("=" + string(make([]byte, 50)) + "=")
			log.Printf("Message received")
			log.Printf("  Topic: %s", message.Topic)
			log.Printf("  Partition: %d", message.Partition)
			log.Printf("  Offset: %d", message.Offset)
			log.Printf("  Key: %s", string(message.Key))
			log.Printf("  Transaction: %+v", transaction)
			log.Println("=" + string(make([]byte, 50)) + "=")
			log.Println()

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			log.Println("Session context done, exiting consume claim")
			return nil
		}
	}
}
