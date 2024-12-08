package app

import (
	"context"
	"producer-simulator/internal/config"
	"producer-simulator/internal/producer"
	"producer-simulator/internal/scheduler"
	"producer-simulator/internal/usecase"

	"github.com/sirupsen/logrus"
)

type App struct {
	sched scheduler.IScheduler
}

func New(ctx context.Context, cfg *config.Config, log *logrus.Logger) *App {
	producer, err := producer.New(cfg.Producer)
	if err != nil {
		panic(err)
	}

	usecases := usecase.NewUseCases(producer)
	sched := scheduler.NewScheduler(usecases.IOrder, log, cfg.Schedular.Period)

	return &App{
		sched: sched,
	}
}

func (a *App) Start(ctx context.Context) {
	a.sched.Start(ctx)
}

func (a *App) Stop() {
	a.sched.Stop()
}
