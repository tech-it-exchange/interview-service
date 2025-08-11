package main

import (
	"context"
	"fmt"
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
	_ "interview-service/docs"
	microserviceConfig "interview-service/infrastructure/config"
	fasthttp2 "interview-service/infrastructure/fasthttp"
	"interview-service/infrastructure/management"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
)

func init() {
	// Грузим локальные env без выхода за пределы проекта
	if err := godotenv.Load(".env.local", ".env"); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := microserviceConfig.GetConfig()

	logger := log.New(os.Stdout, "[app] ", log.LstdFlags)

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
		logger.Printf("Ошибка инициализации потребителей кафки: %v", err)

		return
	}
	err = kafkaAdapter.StartConsuming()
	if err != nil {
		logger.Printf("Ошибка запуска потребителей кафки: %v", err)

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
	logger *log.Logger,
) {
	httpAdapter, err := appFasthttp.NewHttpAdapter(
		controllerDependencies,
	)
	if err != nil {
		logger.Printf("Ошибка инициализации http адаптера: %v", err)

		return
	}
	httpAdapter.InitRoutes()
	routerHandler := httpAdapter.GetHandler()

	port := fmt.Sprintf(":%s", config.HttpServer.HttpPort)
	handler := fasthttp2.NewHandler(logger, routerHandler)

	logger.Printf("Запуск http сервера на порту %s", port)
	err = fasthttp.ListenAndServe(port, handler)
	if err != nil {
		err = fmt.Errorf("не удалось запустить HTTP-сервер, err: %v", err)
		logger.Printf(err.Error())

		panic(err)
	}
}
