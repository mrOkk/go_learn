package kafka

import (
	"App/internal/application"
	"log"

	"github.com/IBM/sarama"
)

type ConsumerHandler struct {
	processor *application.Processor
	queueSize int
}

func NewConsumerHandler(processor *application.Processor, queueSize int) *ConsumerHandler {
	return &ConsumerHandler{
		processor: processor,
		queueSize: queueSize,
	}
}

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
	worker := newPartitionWorker(claim.Partition(), h.processor, h.queueSize)
	defer worker.Close()

	for msg := range claim.Messages() {
		worker.Submit(session, msg)
	}

	return nil
}