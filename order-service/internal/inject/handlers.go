package inject

import (
	"wb-orders/internal/handler"
	"wb-orders/internal/routing"

	"github.com/sirupsen/logrus"
)

func NewHandlers(log *logrus.Logger, usecases UseCases) *routing.Handlers {
	return &routing.Handlers{
		IOrder: handler.NewOrder(log, usecases.IOrder),
	}
}
