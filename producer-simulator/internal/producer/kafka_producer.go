package producer

import (
	"producer-simulator/internal/config"
	"producer-simulator/pkg/generate"

	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type IProducer interface {
	Produce(topic string, payload interface{}) error
}

type KafkaProducer struct {
	p *kafka.Producer
}

func New(cfg config.Producer) (*KafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Servers,
		"security.protocol": cfg.Protocol,
		"acks":              cfg.Acks})

	if err != nil {
		return nil, err
	}
	return &KafkaProducer{
		p: p,
	}, nil
}

func (p *KafkaProducer) Produce(topic string, payload interface{}) error {
	msg, err := encodeMessage(topic, payload)
	if err != nil {
		return err
	}

	p.p.Produce(msg, nil)
	return nil
}

func encodeMessage(topic string, payload interface{}) (*kafka.Message, error) {
	m, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	key := generate.NewUUID().String()

	return &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          []byte(m),
	}, nil
}
