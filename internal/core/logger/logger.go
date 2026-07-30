package core_logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerContextKey struct{}

var (
	key = loggerContextKey{}
)

type Logger struct {
	*zap.Logger //объявляем не как поле, а как встраиваемый тип
	//поэтому можем напрямую использовать методы zap.Logger, например, logger.Info("message"), а не logger.Logger.Info("message")
	file *os.File
}

func ToContext(ctx context.Context, log *Logger) context.Context {
	return context.WithValue(
		ctx,
		key,
		log,
	)
}

// пишем новый конструктор, который из контекста будет доставать логгер, который мы положили туда в middleware Logger(common.go)
func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value(key).(*Logger) //получаем логгер из контекста, который мы положили туда в middleware Logger
	if !ok {
		panic("logger not found in context") //если логгер не найден в контексте, то вызываем панику, чтобы остановить выполнение программы и вывести сообщение об ошибке
	}
	return log
}

func NewLogger(config Config) (*Logger, error) {
	zapLvl := zap.NewAtomicLevel()
	if err := zapLvl.UnmarshalText([]byte(config.Level)); err != nil { //внутрь кладем уровень логирования, который мы хотим использовать, например, zap.InfoLevel
		return nil, fmt.Errorf("invalid log level: %w", err) //Errorf - метод для форматирования строки с ошибкой, %w - используется для обертывания ошибки, чтобы сохранить ее тип и стек вызовов
		//Errorf добавляет контекст к ошибке, что позволяет лучше понимать, что произошло не так, когда ошибка будет обработана в другом месте кода
	}
	//создаем файл для логов
	if err := os.MkdirAll(config.Folder, 0755); err != nil { //0755 - права доступа к папке, 0 - для указания, что это папка, 7 - для владельца (чтение, запись и выполнение), 5 - для группы (чтение и выполнение), 5 - для остальных (чтение и выполнение)
		return nil, fmt.Errorf("mkdir lo")
	}
	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	logFilePath := filepath.Join(
		config.Folder,
		fmt.Sprintf("%s.log", timestamp), //форматируем имя файла с логами, чтобы оно содержало дату и время создания файла
	)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644) //OpenFile - метод для открытия файла, os.O_CREATE - создать файл, если он не существует, os.O_WRONLY - открыть файл только для записи, 0644 - права доступа к файлу, 0 - для указания, что это файл, 6 - для владельца (чтение и запись), 4 - для группы (только чтение), 4 - для остальных (только чтение
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	//создаем конфигурацию для zap.Logger, которая будет писать логи в файл и в консоль
	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvl), //логируем в консоль
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLvl),   //логируем в файл
	)
	zapLogger := zap.New(core, zap.AddCaller())

	return &Logger{
		Logger: zapLogger,
		file:   logFile,
	}, nil
}

// With создает новый логгер с добавленными полями, которые будут автоматически добавляться ко всем логам. Например, если мы хотим добавить request_id и url к каждому логу, мы можем использовать метод With, чтобы создать новый логгер с этими полями.
func (l *Logger) With(field ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(field...),
		file:   l.file,
	}
}

func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Println("failed to close application logger:", err) //в файл не записываем, потому что если он не закрылся, с ним самим проблема скорее всего
	}
}
