package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/Emilia20112005/golang-todoapp/internal/core/config"
	core_logger "github.com/Emilia20112005/golang-todoapp/internal/core/logger"
	core_postgres_pool "github.com/Emilia20112005/golang-todoapp/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/Emilia20112005/golang-todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/Emilia20112005/golang-todoapp/internal/core/transport/http/server"
	statistics_postgres_repository "github.com/Emilia20112005/golang-todoapp/internal/features/statistics/respository/postgres"
	statistics_service "github.com/Emilia20112005/golang-todoapp/internal/features/statistics/service"
	statisctics_transport_http "github.com/Emilia20112005/golang-todoapp/internal/features/statistics/transport/http"
	tasks_postgres_repository "github.com/Emilia20112005/golang-todoapp/internal/features/tasks/repository/postgres"
	tasks_service "github.com/Emilia20112005/golang-todoapp/internal/features/tasks/service"
	tasks_transport_http "github.com/Emilia20112005/golang-todoapp/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/Emilia20112005/golang-todoapp/internal/features/users/repository/postgres"
	users_service "github.com/Emilia20112005/golang-todoapp/internal/features/users/service"
	users_transport_http "github.com/Emilia20112005/golang-todoapp/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext( //контекст для httpServer.Run()
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("Failed to init application logger: ", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("Application time zone", zap.Any("zone", time.Local))

	//создаем пул подключений
	logger.Debug("Initializing postgres connection pool")
	pool, err := core_postgres_pool.NewConnectionPool(
		ctx,
		core_postgres_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	//1. Начинаем выполнение фичи юзерс
	logger.Debug("Initializing feature", zap.String("feature", "users")) //указываем что начинаем выполнение фичи юзерс
	//репозиторий
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	//сервис
	usersService := users_service.NewUsersService(usersRepository)
	//транспорт
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	//2. Начинаем выполнение фичи tasks
	logger.Debug("Initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	//3. Начинаем выполнение фичи statistics
	logger.Debug("Initializing feature", zap.String("feature", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHTTP := statisctics_transport_http.NewStatisticsHTTPHandler(statisticsService)

	//запускаем сервер
	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersionV1)
	apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(statisticsTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouterV1) //здесь по указателю внутрь передаем

	if err := httpServer.Run(ctx); err != nil {
		logger.Error(("HTTP server run error"), zap.Error(err))
	}
}
