package usecase

import (
	"context"
	"producer-simulator/internal/domain/model"
	"producer-simulator/internal/producer"

	"github.com/brianvoe/gofakeit/v6"
)

const (
	ordersTopic = "orders"
)

type IOrder interface {
	Send(ctx context.Context) error
}

type Order struct {
	publisher producer.IProducer
}

func NewOrder(publisher producer.IProducer) *Order {
	return &Order{
		publisher: publisher,
	}
}

func (u *Order) Send(ctx context.Context) error {
	order := generateRandomOrder()

	if err := u.publisher.Produce(ordersTopic, order); err != nil {
		return err
	}

	return nil
}

func generateRandomOrder() model.Order {
	return model.Order{
		OrderUID:    gofakeit.UUID(),
		TrackNumber: gofakeit.LetterN(10),
		Entry:       gofakeit.Word(),
		Delivery: model.Delivery{
			Name:    gofakeit.Name(),
			Phone:   gofakeit.Phone(),
			Zip:     gofakeit.Zip(),
			City:    gofakeit.City(),
			Address: gofakeit.Address().Address,
			Region:  gofakeit.State(),
			Email:   gofakeit.Email(),
		},
		Payment: model.Payment{
			Transaction:  gofakeit.UUID(),
			RequestID:    gofakeit.UUID(),
			Currency:     gofakeit.CurrencyShort(),
			Provider:     gofakeit.Company(),
			Amount:       gofakeit.Number(1000, 50000),
			PaymentDt:    gofakeit.Date().Unix(),
			Bank:         gofakeit.Company(),
			DeliveryCost: gofakeit.Number(100, 1000),
			GoodsTotal:   gofakeit.Number(1000, 50000),
			CustomFee:    gofakeit.Number(0, 500),
		},
		Items:             generateRandomItems(gofakeit.Number(1, 5)),
		Locale:            gofakeit.Language(),
		InternalSignature: gofakeit.UUID(),
		CustomerID:        gofakeit.UUID(),
		DeliveryService:   gofakeit.Company(),
		ShardKey:          gofakeit.Word(),
		SmID:              gofakeit.Number(1, 100),
		DateCreated:       gofakeit.Date(),
		OofShard:          gofakeit.Word(),
	}
}

func generateRandomItems(count int) []model.Item {
	items := make([]model.Item, count)
	for i := 0; i < count; i++ {
		items[i] = model.Item{
			ChrtID:      gofakeit.Number(1000, 9999),
			TrackNumber: gofakeit.LetterN(10),
			Price:       gofakeit.Number(100, 10000),
			RID:         gofakeit.UUID(),
			Name:        gofakeit.ProductName(),
			Sale:        gofakeit.Number(0, 50),
			Size:        gofakeit.LetterN(2),
			TotalPrice:  gofakeit.Number(100, 20000),
			NmID:        gofakeit.Number(100000, 999999),
			Brand:       gofakeit.Company(),
			Status:      gofakeit.Number(0, 5),
		}
	}
	return items
}
