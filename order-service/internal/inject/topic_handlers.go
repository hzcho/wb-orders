package inject

import (
	"wb-orders/internal/listener"
	topichandler "wb-orders/internal/topic_handler"
)

func NewTopicHandlers(usecases UseCases) *listener.Handlers {
	return &listener.Handlers{
		IOrder: topichandler.NewOrder(usecases.IOrder),
	}
}
