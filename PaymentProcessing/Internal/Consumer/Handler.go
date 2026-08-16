package Consumer

import (
	"log"

	"github.com/IBM/sarama"
)

type Handler struct{}

func (h *Handler) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Consumer group session setup")
	return nil
}

func (h *Handler) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("Consumer group session cleanup")
	return nil
}

func (h *Handler) ConsumeClaim(
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

			log.Println("=" + string(make([]byte, 50)) + "=")
			log.Printf("Message received")
			log.Printf("  Topic: %s", message.Topic)
			log.Printf("  Partition: %d", message.Partition)
			log.Printf("  Offset: %d", message.Offset)
			log.Printf("  Key: %s", string(message.Key))
			log.Printf("  Value: %s", string(message.Value))
			log.Println("=" + string(make([]byte, 50)) + "=")
			log.Println()

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			log.Println("Session context done, exiting consume claim")
			return nil
		}
	}
}
