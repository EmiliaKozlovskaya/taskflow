package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_kafka "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/kafka"
	core_logger "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

//Воркер не запускает HTTP, не подключается к PostgreSQL и не содержит tasks service.
// Он делает только:
// подключиться к Kafka -> прочитать TaskEvent -> залогировать -> подтвердить Kafka, что сообщение обработано

func main() {
	//создаем контекст, который будет отменен при получении сигнала SIGINT или SIGTERM
	//для того чтобы корректно завершить работу воркера, когда мы его остановим через docker compose down или через ctrl+c
	//когда контекст отменяется, то все горутины, которые его используют, должны завершить свою работу и освободить ресурсы
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("Failed to init application logger: ", err)
		os.Exit(1)
	}
	defer logger.Close()

	cfg := core_kafka.NewConfigMust()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		Topic:   cfg.Topic,
		GroupID: cfg.GroupID,
	})
	defer reader.Close()

	logger.Info("Kafka worker started")

	for {
		message, err := reader.FetchMessage(ctx)
		if err != nil {
			//ctx.Err() возвращает ошибку, если контекст был отменен, то есть если мы получили сигнал SIGINT или SIGTERM
			if ctx.Err() != nil {
				logger.Info("Kafka worker stopped")
				return
			}

			logger.Error("Failed to fetch message from Kafka", zap.Error(err))
			continue
		}
		var event core_kafka.TaskEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			logger.Error(
				"failed to decode Kafka event",
				zap.Error(err),
				zap.ByteString("message", message.Value), //логируем сообщение, которое не удалось декодировать, чтобы потом можно было его проанализировать
			)
			//подтверждаем Kafka, что сообщение обработано, чтобы оно не повторялось
			if err := reader.CommitMessages(ctx, message); err != nil {
				logger.Error("failed to commit invalid Kafka message", zap.Error(err))
			}
			continue
		}
		logger.Info(
			"task event recieved",
			zap.Time("requested_at", event.Timestamp),
			zap.String("action", string(event.Action)),
		)
		//Подтверждаем сообщение только после записи в лог
		//Разделяем FetchMessage и CommitMessage (вместо ReadMessage которое делает коммит под капотом)
		//Потому что если на этапе после Fetch и до Commit воркер упадет, то после восстановления он вновь обработает
		//данное сообщение, так как оно еще не было зафиксировано
		if err := reader.CommitMessages(ctx, message); err != nil {
			logger.Error("failed to commit Kafka message", zap.Error(err))
		}
	}
}
