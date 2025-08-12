package main

import (
	"context"
	"fmt"
	"github.com/valyala/fasthttp"
	appFasthttp "interview-service/adapter/inbound/fasthttp"
	"interview-service/adapter/inbound/fasthttp/controllers"
	inboundKafka "interview-service/adapter/inbound/kafka"
	"interview-service/adapter/inbound/kafka/manager"
	"interview-service/adapter/outbound/localStorage"
	"interview-service/adapter/outbound/postgres"
	"interview-service/application/channels"
	"interview-service/application/service"
	"interview-service/application/usecase"
	"interview-service/application/worker"
	microserviceConfig "interview-service/infrastructure/config"
	fasthttp2 "interview-service/infrastructure/fasthttp"
	"interview-service/infrastructure/management"
	"log/slog"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := microserviceConfig.GetConfig()

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(handler)

	management.NewHandlerBuilder(logger, config.Management.ManagementPort).
		AddSwaggerUIHandler(management.SwaggerConfig{
			SwaggerUrl:  config.Swagger.SwaggerUrl,
			SwaggerPath: config.Swagger.SwaggerPath,
		}).
		Build()

	storageContainer := localStorage.NewStorageContainer()
	pgRepositories := postgres.NewRepositories()

	channelContainer := channels.NewChannelContainer(config.ChannelConfig, logger)

	services := service.NewServices(
		logger,
		storageContainer,
	)

	useCases := usecase.NewUseCases(
		storageContainer,
		pgRepositories,
		logger,
	)

	workers := worker.NewWorkers(
		storageContainer,
		services,
		logger,
		config.RefreshInstrumentsTimeout,
		channelContainer.WorkerChannelManager,
		channelContainer.InstrumentChannelManager,
	)
	workers.StartWorkers(ctx)

	kafkaManager := manager.NewManager(
		logger,
	)

	kafkaAdapter := inboundKafka.NewKafkaAdapter(
		kafkaManager,
		channelContainer.InstrumentChannelManager,
		services.Kafka,
		logger,
	)

	err := kafkaAdapter.InitConsumers()
	if err != nil {
		logger.Error("Ошибка инициализации потребителей кафки", "err", err)

		return
	}
	err = kafkaAdapter.StartConsuming()
	if err != nil {
		logger.Error("Ошибка запуска потребителей кафки", "err", err)

		return
	}

	defer kafkaAdapter.CloseConsuming()

	controllersDependencies := &controllers.Dependencies{
		UseCases: useCases,
	}

	startHttpServer(
		controllersDependencies,
		config,
		logger,
	)
}

// startHttpServer Запускает http сервер
func startHttpServer(
	controllerDependencies *controllers.Dependencies,
	config *microserviceConfig.MicroserviceConfig,
	logger *slog.Logger,
) {
	httpAdapter, err := appFasthttp.NewHttpAdapter(
		controllerDependencies,
	)
	if err != nil {
		logger.Error("Ошибка инициализации http адаптера", "err", err)

		return
	}
	httpAdapter.InitRoutes()
	routerHandler := httpAdapter.GetHandler()

	port := fmt.Sprintf(":%s", config.HttpServer.HttpPort)
	handler := fasthttp2.NewHandler(logger, routerHandler)

	logger.Info("Запуск http сервера", "port", port)
	err = fasthttp.ListenAndServe(port, handler)
	if err != nil {
		err = fmt.Errorf("не удалось запустить HTTP-сервер, err: %v", err)
		logger.Error(err.Error(), "err", err)

		panic(err)
	}
}
