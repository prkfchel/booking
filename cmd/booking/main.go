package main

import (
	core_config "booking/internal/core/config"
	core_logger "booking/internal/core/logger"
	core_pgx_pool "booking/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "booking/internal/core/transport/http/middlerware"
	core_http_server "booking/internal/core/transport/http/server"
	bookings_postgres_repository "booking/internal/features/bookings/repository/postgres"
	bookings_service "booking/internal/features/bookings/service"
	bookings_http_transport "booking/internal/features/bookings/transport/http"
	events_postgres_repository "booking/internal/features/events/repository/postgres"
	events_service "booking/internal/features/events/service"
	events_http_transport "booking/internal/features/events/transport/http"
	users_postgres_repository "booking/internal/features/users/repository/postgres"
	users_service "booking/internal/features/users/service"
	users_http_transport "booking/internal/features/users/transport/http"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "booking/docs"

	"go.uber.org/zap"
)

// @title 			Booking API
// @version 		1.0
// @description 	Booking App REST-API scheme
// @host 			127.0.0.1:5050
// @BasePath 		/api/v1

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	pool, err := core_pgx_pool.NewPool(ctx, core_pgx_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersHTTPHandler := users_http_transport.NewHTTPUsersHandler(usersService)

	eventsRepository := events_postgres_repository.NewEventsRepository(pool)
	eventsServise := events_service.NewEventsService(eventsRepository)
	eventsHTTPHandler := events_http_transport.NewEventsHTTPHandler(eventsServise)

	bookRepository := bookings_postgres_repository.NewBookRepository(pool)
	bookService := bookings_service.NewBookService(bookRepository)
	bookHTTPHandler := bookings_http_transport.NewHTTPBookHandler(bookService)

	httpConfig := core_http_server.NewConfigMust()

	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(usersHTTPHandler.Routes()...)
	apiVersionRouterV1.RegisterRoutes(eventsHTTPHandler.Routes()...)
	apiVersionRouterV1.RegisterRoutes(bookHTTPHandler.Routes()...)

	httpServer.RegisterAPIRoutes(
		apiVersionRouterV1,
	)

	httpServer.RegisterSwagger()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("server error", zap.Error(err))
	}
}
