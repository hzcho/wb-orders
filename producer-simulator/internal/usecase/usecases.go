package usecase

import "producer-simulator/internal/producer"

type UseCases struct {
	IOrder
}

func NewUseCases(publisher producer.IProducer) *UseCases {
	return &UseCases{
		IOrder: NewOrder(publisher),
	}
}
