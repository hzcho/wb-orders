package listener

import (
	"context"
	"wb-orders/internal/consumer"
	topichandler "wb-orders/internal/topic_handler"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	topichandler.IOrder
}

type TopicFunc func(context.Context, *kafka.Message) error

type KafkaListener struct {
	consumer      consumer.IConsumer
	log           *logrus.Logger
	stop          chan struct{}
	topicHandlers map[string]TopicFunc
}

func New(consumer consumer.IConsumer, log *logrus.Logger) *KafkaListener {
	stop := make(chan struct{})
	topicHandlers := make(map[string]TopicFunc)

	//topicHandlers["save_operation"] = handlers.Operation.Save

	return &KafkaListener{
		consumer:      consumer,
		log:           log,
		stop:          stop,
		topicHandlers: topicHandlers,
	}
}

func (k *KafkaListener) SetHandlers(handlers Handlers) {
	k.topicHandlers["orders"] = handlers.IOrder.Save
}

func (k *KafkaListener) Start(ctx context.Context) {
	for {
		select {
		case <-k.stop:
			return
		default:
			ms, err := k.consumer.Read()
			if err != nil {
				if err.Error() != "Local: Timed out" {
					k.log.Error(err)
				}
				continue
			}

			go func(ms *kafka.Message) {
				err := k.topicHandlers[*ms.TopicPartition.Topic](ctx, ms)
				if err != nil {
					k.log.Error(err)
				}
			}(ms)
		}
	}
}

func (k *KafkaListener) Stop() {
	k.stop <- struct{}{}
	defer close(k.stop)
}
