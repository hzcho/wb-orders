package usecase

import (
	"context"
	"fmt"
	"wb-orders/internal/cache"
	"wb-orders/internal/domain/model"
	"wb-orders/internal/domain/request"
	"wb-orders/internal/domain/response"
	"wb-orders/internal/repository"
)

type IOrder interface {
	GetOrderIDs(ctx context.Context, req request.Pagination) (response.GetOrderIDs, error)
	Save(ctx context.Context, order model.Order) error
	GetById(ctx context.Context, id string) (model.Order, error)
	LoadCache(ctx context.Context, limit int) error
}

type Order struct {
	orderRepo repository.IOrder
	cache     cache.ICache
}

func NewOrder(orderRepo repository.IOrder, cache cache.ICache) *Order {
	return &Order{
		orderRepo: orderRepo,
		cache:     cache,
	}
}

func (u *Order) GetOrderIDs(ctx context.Context, req request.Pagination) (response.GetOrderIDs, error) {

	if req.Page < 0 {
		req.Page = 0
	}
	if req.Limit < 1 {
		req.Limit = 10
	}

	orders, err := u.orderRepo.GetAll(ctx, req.Limit, req.Page)
	if err != nil {
		return response.GetOrderIDs{}, err
	}

	ids := make([]string, 0, len(orders))

	for _, v := range orders {
		ids = append(ids, v.OrderUID)
	}

	return response.GetOrderIDs{
		Page:  req.Page,
		Limit: req.Limit,
		IDs:   ids,
	}, nil
}

func (u *Order) Save(ctx context.Context, order model.Order) error {
	err := u.orderRepo.Save(ctx, order)
	if err != nil {
		return err
	}

	u.cache.Set(order.OrderUID, order)

	return nil
}

func (u *Order) GetById(ctx context.Context, id string) (model.Order, error) {
	if order, exists := u.cache.Get(id); exists {
		return order, nil
	}

	order, err := u.orderRepo.GetById(ctx, id)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (u *Order) LoadCache(ctx context.Context, limit int) error {
	orders, err := u.orderRepo.GetAll(ctx, limit, 0)
	if err != nil {
		return fmt.Errorf("failed to load orders from database: %w", err)
	}

	items := make(map[string]model.Order, len(orders))
	for _, order := range orders {
		items[order.OrderUID] = order
	}

	u.cache.SetMany(items)
	return nil
}
