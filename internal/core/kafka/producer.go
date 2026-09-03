package core_kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

//здесь реальная реализация интерфейса TaskEventPublisher, которая будет публиковать события в Kafka. Внутри будет использоваться библиотека segmentio/kafka-go,
// которая предоставляет удобный API для работы с Kafka. Внутри будет создан Kafka Writer, который будет отправлять сообщения в указанный топик.
// В конструкторе NewKafkaPublisher мы будем принимать конфигурацию Kafka, чтобы знать, куда отправлять сообщения.

type KafkaPublisher struct {
	writer *kafka.Writer
}

// передаем конфиг, чтобы знать, куда отправлять сообщения. Внутри создаем Kafka Writer, который будет отправлять сообщения в указанный топик.
func NewKafkaPublisher(cfg Config) *KafkaPublisher {
	return &KafkaPublisher{
		writer: &kafka.Writer{
			Addr:  kafka.TCP(cfg.Brokers...),
			Topic: cfg.Topic,
		},
	}
}

func (p *KafkaPublisher) Publish(
	ctx context.Context,
	event TaskEvent,
) error {
	message, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal task event: %w", err)
	}

	if err := p.writer.WriteMessages(ctx, kafka.Message{ //физически отправляет JSON в topic todo.task-events
		Value: message,
	}); err != nil {
		return fmt.Errorf("write Kafka message: %w", err)
	}

	return nil
}

func (p *KafkaPublisher) Close() error {
	return p.writer.Close()
}
