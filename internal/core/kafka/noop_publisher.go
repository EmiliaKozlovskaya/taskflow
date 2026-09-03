package core_kafka

import "context"

// в случане если Kafka отключена, то мы можем использовать NoopPublisher, который ничего не делает, но при этом не ломает работу приложения. Это удобно для локальной разработки и тестирования, когда Kafka может быть недоступна или не нужна.
// KAFKA_ENABLED=false
type NoopPublisher struct{}

func (NoopPublisher) Publish(
	ctx context.Context,
	taskEvent TaskEvent,
) error {
	return nil
}

func (NoopPublisher) Close() error {
	return nil
}
