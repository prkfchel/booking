package main

import (
	core_config "booking/internal/core/config"
	core_logger "booking/internal/core/logger"
	core_pgx_pool "booking/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "booking/internal/core/transport/http/middlerware"
	core_http_server "booking/internal/core/transport/http/server"
	users_postgres_repository "booking/internal/features/users/repository/postgres"
	users_service "booking/internal/features/users/service"
	users_http_transport "booking/internal/features/users/transport/http"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

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
	usersTransport := users_http_transport.NewHTTPUsersHandler(usersService)

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
	apiVersionRouterV1.RegisterRoutes(usersTransport.Routes()...)

	httpServer.RegisterAPIRoutes(
		apiVersionRouterV1,
	)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("server error", zap.Error(err))
	}
}
