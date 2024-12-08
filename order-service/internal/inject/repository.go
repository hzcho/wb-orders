package inject

import (
	"wb-orders/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	repository.IOrder
}

func NewRepositories(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		IOrder: repository.NewOrder(pool),
	}
}
