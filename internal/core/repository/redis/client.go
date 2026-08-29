package core_redis

import (
	"context"
	"fmt"
	"net"

	"github.com/redis/go-redis/v9"
)

// Создает нового клиента Redis с использованием предоставленной конфигурации.
// Потом можно будет использовать этот клиент для взаимодействия с Redis.
func NewClient(ctx context.Context, cfg Config) (*redis.Client, error) {
	//redis.Options - структура, которая содержит параметры подключения к Redis.
	client := redis.NewClient(&redis.Options{
		Addr:         net.JoinHostPort(cfg.Host, cfg.Port),
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.Timeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	// Проверяем соединение с Redis, отправляя команду PING.
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close() //намеренно игнорируем ошибку закрытия клиента, так как мы уже получили ошибку при попытке подключения.
		return nil, fmt.Errorf("failed to connect to redis server: %w", err)
	}

	return client, nil
}
