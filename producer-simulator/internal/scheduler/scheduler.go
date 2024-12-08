package scheduler

import (
	"context"
	"time"

	"producer-simulator/internal/usecase"

	"github.com/sirupsen/logrus"
)

type IScheduler interface {
	Start(ctx context.Context)
	Stop()
}

type Scheduler struct {
	log          *logrus.Logger
	stop         chan struct{}
	orderUseCase usecase.IOrder
	period       time.Duration
}

func NewScheduler(orderUseCase usecase.IOrder, log *logrus.Logger, period time.Duration) *Scheduler {
	return &Scheduler{
		log:          log,
		stop:         make(chan struct{}),
		orderUseCase: orderUseCase,
		period:       period,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	go func() {
		timer := time.NewTimer(time.Second * 10)
		defer timer.Stop()

		for {
			select {
			case <-s.stop:
				return
			case <-timer.C:
				if err := s.orderUseCase.Send(ctx); err != nil {
					s.log.WithField("op", "/internal/scheduler/Start").Error(err)
				}
				timer.Reset(s.period)
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.stop)
}
