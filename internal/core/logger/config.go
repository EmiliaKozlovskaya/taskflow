package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

//config.go - файл для конфигурации логгера, например, для указания уровня логирования и папки для логов.
//чтобы параметры в конструкторе func NewLogger(logLevel string, logFolder string) были более понятными
// и не нужно было гадать, что за строку нужно передать в конструктор, например, "info" или "debug", мы можем создать константы для уровней логирования.

type Config struct {
	Level  string `envconfig:"LEVEL" default:"DEBUG"` //envconfig - библиотека, которая позволяет загружать значения из переменных окружения в структуру. required: "true" - поле обязательно.
	Folder string `envconfig:"FOLDER" required:"true"`
} //можно задавать разными способами, например, через переменные окружения (самый популярный), конфигурационные файлы или аргументы командной строки.

// конструктор для LoggerConfig, который будет загружать конфигурацию из переменных окружения.
func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("LOGGER", &config); err != nil { //Process - метод, который загружает значения из переменных окружения, которые начинаются с "LOGGER", в структуру config. Например, LOGGER_LEVEL и LOGGER_FOLDER.
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}
	return config, nil
}

// Этот конструктор с Must вызывается, когда конфиг логгера ДОЛЖЕН появиться, должен корректно быть прочитан, а если этого не происходит, то кидаем панику
// Нужен чтобы не возиться лишний раз с обработкой ошибок, а сразу увидеть в чем проблема.
func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get Logger config: %w", err)
		panic(err)
	}
	return config
}
