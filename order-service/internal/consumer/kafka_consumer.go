package consumer

import (
	"time"
	"wb-orders/internal/config"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type IConsumer interface {
	Read() (*kafka.Message, error)
}

type KafkaConsumer struct {
	c *kafka.Consumer
}

func NewKafkaConsumer(cfg config.Consumer) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Brokers,
		"group.id":          cfg.GroupId,
		"auto.offset.reset": cfg.Offset,
	})
	if err != nil {
		return nil, err
	}

	c.SubscribeTopics(cfg.Topics, nil)

	return &KafkaConsumer{
		c: c,
	}, nil
}

func (c *KafkaConsumer) Read() (*kafka.Message, error) {
	ev, err := c.c.ReadMessage(100 * time.Millisecond)
	if err != nil {
		return nil, err
	}

	return ev, nil
}
