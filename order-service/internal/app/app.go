package app

import (
	"context"
	"fmt"
	"wb-orders/internal/cache"
	"wb-orders/internal/config"
	"wb-orders/internal/consumer"
	"wb-orders/internal/inject"
	"wb-orders/internal/listener"
	"wb-orders/internal/routing"
	"wb-orders/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type App struct {
	pool     *pgxpool.Pool
	server   *server.Server
	listener *listener.KafkaListener
}

func New(ctx context.Context, cfg *config.Config, log *logrus.Logger) (*App, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.PG.Username, cfg.PG.Password, cfg.PG.Host, cfg.PG.Port, cfg.PG.DBName)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database connection pool: %w", err)
	}

	c, err := consumer.NewKafkaConsumer(cfg.Consumer)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to initialize Kafka consumer: %w", err)
	}

	cache := cache.NewCache(cfg.Cache.Capacity)
	repos := inject.NewRepositories(pool)
	usecases := inject.NewUseCases(*repos, cache)
	handlers := inject.NewHandlers(log, *usecases)
	topicHandlers := inject.NewTopicHandlers(*usecases)

	if err := usecases.IOrder.LoadCache(ctx, cfg.Cache.Capacity); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to load cache: %w", err)
	}

	lstnr := listener.New(c, log)
	lstnr.SetHandlers(*topicHandlers)

	router := gin.New()
	routing.InitRoutes(router, *handlers)

	server := server.New(&cfg.Server, router)

	return &App{
		pool:     pool,
		server:   server,
		listener: lstnr,
	}, nil
}

func (a *App) Start(ctx context.Context) {
	go a.server.Run()
	go a.listener.Start(ctx)
}

func (a *App) Stop(ctx context.Context) {
	a.pool.Close()
	a.server.Stop(ctx)
	a.listener.Stop()
}
