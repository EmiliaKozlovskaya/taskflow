package core_kafka

import (
	"context"
	"time"
)

type TaskAction string

const (
	TaskCreated TaskAction = "task_created"
	TaskUpdated TaskAction = "task_updated"
	TaskDeleted TaskAction = "task_deleted"
	TasksListed TaskAction = "tasks_listed"
	TaskListed  TaskAction = "task_listed"
)

type TaskEvent struct {
	Timestamp time.Time  `json:"timestamp"`
	Action    TaskAction `json:"action"`
}

func NewTaskEvent(
	timestamp time.Time,
	action TaskAction,
) TaskEvent {
	return TaskEvent{
		Timestamp: timestamp,
		Action:    action,
	}
}

type TaskEventPublisher interface {
	Publish(ctx context.Context, event TaskEvent) error
	Close() error //нужен потому что реальный Kafka Producer в конце работы должен корректно закрыться, чтобы не было утечек памяти и не потерялись сообщения, которые еще не успели уйти в брокер.
}
