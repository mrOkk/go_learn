package kafka

import (
	"App/internal/application"
	"App/internal/domain"
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type workItem struct {
	session sarama.ConsumerGroupSession
	message *sarama.ConsumerMessage
}

type partitionWorker struct {
	partition int32
	processor *application.Processor

	jobs chan workItem
	done chan struct{}
}

func newPartitionWorker(partition int32, processor *application.Processor, queueSize int) *partitionWorker {
	w := &partitionWorker{
		partition: partition,
		processor: processor,
		jobs:      make(chan workItem, queueSize),
		done:      make(chan struct{}),
	}

	go w.run()
	return w
}

func (p *partitionWorker) Submit(session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage) {
	p.jobs <- workItem{
		session: session,
		message: message,
	}
}

func (p *partitionWorker) Close() {
	close(p.jobs)
	<-p.done
}

func (p *partitionWorker) run() {
	defer close(p.done)

	for item := range p.jobs {
		var tx domain.Transaction
		if err := json.Unmarshal(item.message.Value, &tx); err != nil {
			log.Printf("[partition=%d offset=%d] decode error: %v", p.partition, item.message.Offset, err)
			continue
		}

		if err := p.processor.ProcessTransaction(context.Background(), tx); err != nil {
			log.Printf("[partition=%d offset=%d] process error: %v", p.partition, item.message.Offset, err)
			continue
		}

		item.session.MarkMessage(item.message, "")
		log.Printf("[partition=%d offset=%d] processed", p.partition, item.message.Offset)
	}
}
