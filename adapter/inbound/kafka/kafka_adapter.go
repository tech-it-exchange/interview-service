package kafka

import (
	"context"
	"interview-service/adapter/inbound/kafka/consumer"
	"interview-service/adapter/inbound/kafka/handlers"
	"interview-service/adapter/inbound/kafka/manager"
	"interview-service/application/channels"
	"interview-service/application/service"
	kafkaPort "interview-service/port/inbound/kafka"
	"log"

	"github.com/google/uuid"
)

type KafkaAdapter struct {
	strategyManager          manager.StrategyManagerInterface
	logger                   *log.Logger
	instrumentChannelManager *channels.InstrumentChannelManager
	kafkaService             *service.KafkaService
}

func NewKafkaAdapter(
	strategyManager manager.StrategyManagerInterface,
	instrumentChannelManager *channels.InstrumentChannelManager,
	kafkaService *service.KafkaService,
	logger *log.Logger,
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
			a.logger.Printf("Не удалось создать потребителя кафки topic=%s err=%v", topic, err)

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

		a.logger.Printf("Ошибка закрытия подключения к кафке по топику=%s err=%v", topic, err)
	}
}

// ListenForNewInstrument Слушает спот инструменты и запускает новые потребители по мере их поступления
func (a *KafkaAdapter) ListenForNewInstrument(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			a.logger.Printf("Остановка ListenForNewInstrument: %v", ctx.Err())
			return
		case InstrumentId, ok := <-a.instrumentChannelManager.GetKafkaSpotInstrumentIdChannel():
			if !ok {
				a.logger.Printf("Пропуск, канал новых спот инструментов закрыт")
				continue
			}

			topic := a.kafkaService.CreateSpotInstrumentTopic(InstrumentId)
			ok = a.kafkaService.HasKafkaConsumer(topic)
			if ok {
				a.logger.Printf("Потребитель для спот инструмента уже существует: %s", InstrumentId)

				continue
			}

			kafkaConsumer, err := a.initConsumer(topic)
			if err != nil {
				a.logger.Printf("Ошибка инициализации нового потребителя спот инструмента %s err=%v", InstrumentId, err)

				continue
			}

			err = kafkaConsumer.StartConsuming()
			if err != nil {
				a.logger.Printf("Ошибка запуска нового потребителя спот инструмента %s err=%v", InstrumentId, err)
				continue
			}

			a.kafkaService.SaveKafkaConsumer(topic, kafkaConsumer)
			a.kafkaService.SaveActiveSpotInstrument(InstrumentId)

			a.logger.Printf("Запущен новый потребитель для спот инструмента: %s", InstrumentId)
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
		a.logger.Printf("Ошибка создания обработчика сообщений из кафки: %v", err)

		return nil, err
	}

	return kafkaConsumer, nil
}
