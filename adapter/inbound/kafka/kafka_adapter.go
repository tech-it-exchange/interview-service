package kafka

import (
	"context"
	"github.com/google/uuid"
	"interview-service/adapter/inbound/kafka/consumer"
	"interview-service/adapter/inbound/kafka/handlers"
	"interview-service/adapter/inbound/kafka/manager"
	"interview-service/application/channels"
	"interview-service/application/service"
	kafkaPort "interview-service/port/inbound/kafka"
	"log/slog"
)

type KafkaAdapter struct {
	strategyManager          manager.StrategyManagerInterface
	logger                   *slog.Logger
	instrumentChannelManager *channels.InstrumentChannelManager
	kafkaService             *service.KafkaService
}

func NewKafkaAdapter(
	strategyManager manager.StrategyManagerInterface,
	instrumentChannelManager *channels.InstrumentChannelManager,
	kafkaService *service.KafkaService,
	logger *slog.Logger,
) kafkaPort.KafkaAdapterInterface {
	return &KafkaAdapter{
		strategyManager:          strategyManager,
		instrumentChannelManager: instrumentChannelManager,
		kafkaService:             kafkaService,
		logger:                   logger,
	}
}

// InitConsumers Создает потребителей
func (a *KafkaAdapter) InitConsumers() error {
	kafkaConsumersMap := make(map[string]*consumer.KafkaConsumer)
	activeKafkaConsumersMap := make(map[uuid.UUID]bool)

	commandTopics := a.kafkaService.GetCommandTopics()
	confirmTopics := a.kafkaService.GetConfirmTopics()
	topics := append(commandTopics, confirmTopics...)

	for _, topic := range topics {
		kafkaConsumer, err := a.initConsumer(topic)
		if err != nil {
			return err
		}

		kafkaConsumersMap[topic] = kafkaConsumer
	}

	for instrumentId, topic := range a.kafkaService.GetSpotInstrumentTopics() {
		kafkaConsumer, err := a.initConsumer(topic)
		if err != nil {
			return err
		}

		kafkaConsumersMap[topic] = kafkaConsumer
		activeKafkaConsumersMap[instrumentId] = true
	}

	a.kafkaService.SaveKafkaConsumersMap(kafkaConsumersMap)
	a.kafkaService.SaveActiveSpotInstrumentsMap(activeKafkaConsumersMap)
	a.kafkaService.SetSpotNeedInstrumentDispatch(true)

	return nil
}

// StartConsuming Запускает потребителей кафки
func (a *KafkaAdapter) StartConsuming() error {
	for topic, kafkaConsumer := range a.kafkaService.GetKafkaConsumers() {
		err := kafkaConsumer.StartConsuming()
		if err != nil {
			a.logger.Error("Не удалось создать потребителя кафки", "topic", topic, "err", err)

			return err
		}
	}

	return nil
}

// CloseConsuming Закрывает подключение у потребителей кафки
func (a *KafkaAdapter) CloseConsuming() {
	for topic, kafkaConsumer := range a.kafkaService.GetKafkaConsumers() {
		err := kafkaConsumer.Close()
		if err == nil {
			continue
		}

		a.logger.Error("Ошибка закрытия подключения к кафке по топику", "topic", topic, "err", err)
	}
}

// ListenForNewInstrument Слушает спот инструменты и запускает новые потребители по мере их поступления
func (a *KafkaAdapter) ListenForNewInstrument(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			a.logger.Info("Остановка ListenForNewInstrument:", ctx.Err())
			return
		case InstrumentId, ok := <-a.instrumentChannelManager.GetKafkaSpotInstrumentIdChannel():
			if !ok {
				a.logger.Error("Пропуск, канал новых спот инструментов закрыт")
				continue
			}

			topic := a.kafkaService.CreateSpotInstrumentTopic(InstrumentId)
			ok = a.kafkaService.HasKafkaConsumer(topic)
			if ok {
				a.logger.Info("Потребитель для спот инструмента уже существует", "InstrumentId", InstrumentId)

				continue
			}

			kafkaConsumer, err := a.initConsumer(topic)
			if err != nil {
				a.logger.Info("Ошибка инициализации нового потребителя спот инструмента", "InstrumentId", InstrumentId, "err", err)

				continue
			}

			err = kafkaConsumer.StartConsuming()
			if err != nil {
				a.logger.Info("Ошибка запуска нового потребителя спот инструмента", "InstrumentId", InstrumentId, "err", err)
				continue
			}

			a.kafkaService.SaveKafkaConsumer(topic, kafkaConsumer)
			a.kafkaService.SaveActiveSpotInstrument(InstrumentId)

			a.logger.Info("Запущен новый потребитель для спот инструмента", "InstrumentId", InstrumentId)
		}
	}
}

// initConsumer Создает потребителя
func (a *KafkaAdapter) initConsumer(topic string) (*consumer.KafkaConsumer, error) {
	handler := handlers.NewCommonHandler(
		a.strategyManager,
		topic,
		"HandleMarkMessageModeCommitSuccessOnly",
		a.logger,
	)

	kafkaConsumer, err := consumer.NewKafkaConsumer(handler)

	if err != nil {
		a.logger.Error("Ошибка создания обработчика сообщений из кафки", "err", err)

		return nil, err
	}

	return kafkaConsumer, nil
}
