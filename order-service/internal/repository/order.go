package repository

import (
	"context"
	"wb-orders/internal/domain/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IOrder interface {
	GetAll(ctx context.Context, limit, offset int) ([]model.Order, error)
	Save(ctx context.Context, order model.Order) error
	GetById(ctx context.Context, id string) (model.Order, error)
}

type Order struct {
	pool *pgxpool.Pool
}

func NewOrder(pool *pgxpool.Pool) *Order {
	return &Order{
		pool: pool,
	}
}

func (r *Order) GetAll(ctx context.Context, limit, offset int) ([]model.Order, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT 
			order_uid, track_number, entry, locale, internal_signature, 
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders 
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order

	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
			&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
			&order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard,
		)
		if err != nil {
			return nil, err
		}

		err = r.pool.QueryRow(ctx, `
			SELECT 
				name, phone, zip, city, address, region, email
			FROM deliveries 
			WHERE order_uid = $1
		`, order.OrderUID).Scan(
			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
			&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region,
			&order.Delivery.Email,
		)
		if err != nil {
			return nil, err
		}

		err = r.pool.QueryRow(ctx, `
			SELECT 
				transaction, request_id, currency, provider, amount, payment_dt, 
				bank, delivery_cost, goods_total, custom_fee
			FROM payments 
			WHERE order_uid = $1
		`, order.OrderUID).Scan(
			&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
			&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt,
			&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal,
			&order.Payment.CustomFee,
		)
		if err != nil {
			return nil, err
		}

		rowsItems, err := r.pool.Query(ctx, `
			SELECT 
				chrt_id, track_number, price, rid, name, sale, size, 
				total_price, nm_id, brand, status
			FROM items 
			WHERE order_uid = $1
		`, order.OrderUID)
		if err != nil {
			return nil, err
		}

		for rowsItems.Next() {
			var item model.Item
			err := rowsItems.Scan(
				&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name,
				&item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand,
				&item.Status,
			)
			if err != nil {
				return nil, err
			}
			order.Items = append(order.Items, item)
		}
		if rowsItems.Err() != nil {
			return nil, rowsItems.Err()
		}

		orders = append(orders, order)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return orders, nil
}

func (r *Order) Save(ctx context.Context, order model.Order) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO orders (
			order_uid, track_number, entry, locale, internal_signature, 
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO deliveries (
			order_uid, name, phone, zip, city, address, region, email
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO payments (
			order_uid, transaction, request_id, currency, provider, 
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt,
		order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err = tx.Exec(ctx, `
			INSERT INTO items (
				order_uid, chrt_id, track_number, price, rid, name, 
				sale, size, total_price, nm_id, brand, status
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`, order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *Order) GetById(ctx context.Context, id string) (model.Order, error) {
	var order model.Order

	err := r.pool.QueryRow(ctx, `
		SELECT 
			order_uid, track_number, entry, locale, internal_signature, 
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders 
		WHERE order_uid = $1
	`, id).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
		&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
		&order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.Order{}, nil
		}
		return model.Order{}, err
	}

	err = r.pool.QueryRow(ctx, `
		SELECT 
			name, phone, zip, city, address, region, email
		FROM deliveries 
		WHERE order_uid = $1
	`, order.OrderUID).Scan(
		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
		&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region,
		&order.Delivery.Email,
	)
	if err != nil {
		return model.Order{}, err
	}

	err = r.pool.QueryRow(ctx, `
		SELECT 
			transaction, request_id, currency, provider, amount, payment_dt, 
			bank, delivery_cost, goods_total, custom_fee
		FROM payments 
		WHERE order_uid = $1
	`, order.OrderUID).Scan(
		&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
		&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt,
		&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
	)
	if err != nil {
		return model.Order{}, err
	}

	rows, err := r.pool.Query(ctx, `
		SELECT 
			chrt_id, track_number, price, rid, name, sale, size, 
			total_price, nm_id, brand, status
		FROM items 
		WHERE order_uid = $1
	`, order.OrderUID)
	if err != nil {
		return model.Order{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Item
		err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name,
			&item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand,
			&item.Status,
		)
		if err != nil {
			return model.Order{}, err
		}
		order.Items = append(order.Items, item)
	}
	if rows.Err() != nil {
		return model.Order{}, rows.Err()
	}

	return order, nil
}
