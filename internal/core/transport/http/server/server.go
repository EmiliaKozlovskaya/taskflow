package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/Emilia20112005/golang-todoapp/internal/core/logger"
	core_http_middleware "github.com/Emilia20112005/golang-todoapp/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux    *http.ServeMux //мультиплексор - сущность, которая будет обрабатывать http запросы и направлять их к нужным middleware и хэндлерам
	config Config
	log    *core_logger.Logger

	middleware []core_http_middleware.Middleware
}

func NewHTTPServer(
	config Config,
	log *core_logger.Logger,
	middleware ...core_http_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (s *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) { //указатель чтобы не копировать встроенный итак внутрь указатель на мукс
	for _, router := range routers {
		//и теперь регистрируем в глобальном мультиплексоре сервера локальные мультиплексоры каждого отдельного APIVersionRouter
		prefix := "/api/" + string(router.apiVersion)

		s.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router.WithMidlleware()), //тут можем использовать router в качестве http.Handler потому что APIVersionRouter удовлетворяет интерфейсу Handler
		)
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddleware(s.mux, s.middleware...)

	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: mux,
	}
	ch := make(chan error, 1)

	//запускаем сервер в отдельной горутине, чтобы не блокировать основной поток выполнения программы
	go func() {
		defer close(ch)

		s.log.Warn("starting HTTP server", zap.String("addr", s.config.Addr))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) { //если ошибка не связана с закрытием сервера, то отправляем ее в канал ch
			ch <- err
		}
	}()

	select {
	case err := <-ch: //если в канал ch пришла ошибка, то возвращаем ее
		if err != nil {
			return fmt.Errorf("listen and serve HTTP: %w", err)
		}
	case <-ctx.Done(): //если контекст завершился, то останавливаем сервер
		s.log.Warn("shutting down HTTP server...")

		//создаем независимый контекст с таймаутом, чтобы сервер успел корректно завершить все текущие запросы и закрыть соединения
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.config.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil { //просим сервер завершиться аккуратно(сервер перестает принимать новые запросы, но продолжает обрабатывать текущие, пока они не завершатся или не истечет таймаут)
			_ = server.Close() //если не удалось корректно завершить серве, то закрываем его принудительно

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		s.log.Warn("HTTP server stopped")
	}
	return nil
}
