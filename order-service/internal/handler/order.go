package handler

import (
	"net/http"
	"wb-orders/internal/domain/request"
	"wb-orders/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type IOrder interface {
	GetOrderIDs(c *gin.Context)
	GetById(c *gin.Context)
}

type Order struct {
	log          *logrus.Logger
	orderUseCase usecase.IOrder
}

func NewOrder(log *logrus.Logger, orderUseCase usecase.IOrder) *Order {
	return &Order{
		log:          log,
		orderUseCase: orderUseCase,
	}
}

func (h *Order) GetOrderIDs(c *gin.Context) {
	log := h.log.WithFields(logrus.Fields{
		"op": "/internal/handler/order/GetOrderIDs",
	})

	var pagination request.Pagination
	if err := c.ShouldBindQuery(&pagination); err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	orderIDs, err := h.orderUseCase.GetOrderIDs(c, pagination)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.HTML(http.StatusOK, "order_ids.html", orderIDs)
}

func (h *Order) GetById(c *gin.Context) {
	log := h.log.WithFields(logrus.Fields{
		"op": "/internal/handler/order/GetById",
	})

	id := c.Param("id")
	if id == "" {
		log.Error("empty ID")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	order, err := h.orderUseCase.GetById(c, id)
	if err != nil {
		log.WithField("id", id).Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.HTML(200, "order.html", gin.H{
		"order": order,
	})
}
