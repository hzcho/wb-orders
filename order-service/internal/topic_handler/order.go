package topichandler

import (
	"context"
	"encoding/json"
	"errors"
	"wb-orders/internal/domain/model"
	"wb-orders/internal/usecase"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type IOrder interface {
	Save(ctx context.Context, msg *kafka.Message) error
}

type Order struct {
	orderUseCase usecase.IOrder
}

func NewOrder(orderUseCase usecase.IOrder) *Order {
	return &Order{
		orderUseCase: orderUseCase,
	}
}

func (h *Order) Save(ctx context.Context, msg *kafka.Message) error {
	var order model.Order

	if msg == nil {
		return errors.New("kafka message is nil")
	}
	if err := json.Unmarshal(msg.Value, &order); err != nil {
		return err
	}

	err := h.orderUseCase.Save(ctx, order)
	if err != nil {
		return err
	}

	return nil
}
