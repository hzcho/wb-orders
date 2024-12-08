package routing

import (
	"wb-orders/internal/handler"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	handler.IOrder
}

func InitRoutes(router *gin.Engine, handlers Handlers) {

	router.Static("/internal/static", "./static")
	router.LoadHTMLGlob("internal/templates/*")

	api := router.Group("/api")
	{
		orders := api.Group("orders")
		{
			orders.GET("/", handlers.IOrder.GetOrderIDs)
			orders.GET("/:id", handlers.IOrder.GetById)
		}
	}
}
