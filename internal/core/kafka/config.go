package core_kafka

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Enabled и Topic в .env, потому что их чаще возможно нужно менять и проще в enabled, в docker-compose оставляем постоянный внутренний адрес 9092
type Config struct {
	Enabled bool     `envconfig:"ENABLED" default:"false"`
	Brokers []string `envconfig:"BROKERS" default:"localhost:9092"` //массив, потому что в настоящем Kafka-кластере брокеров несколько, но у нас пока только один
	Topic   string   `envconfig:"TOPIC" default:"todo.task-events"`
	GroupID string   `envconfig:"GROUP_ID" default:"todo-kafka-worker"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("KAFKA", &config); err != nil { //заполни поля структуры, начинающиеся с KAFKA_
		return Config{}, fmt.Errorf("process Kafka config: %w", err)
	}
	return config, nil
}

func NewConfigMust() Config {
	cfg, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get Kafka config: %w", err)
		panic(err)
	}
	return cfg
}
