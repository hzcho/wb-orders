package inject

import (
	"wb-orders/internal/cache"
	"wb-orders/internal/usecase"
)

type UseCases struct {
	usecase.IOrder
}

func NewUseCases(repos Repositories, cache cache.ICache) *UseCases {
	return &UseCases{
		IOrder: usecase.NewOrder(repos.IOrder, cache),
	}
}
