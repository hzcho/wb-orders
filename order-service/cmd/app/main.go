package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"wb-orders/internal/app"
	"wb-orders/internal/config"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	log := logrus.New()
	cfg, err := config.InitConfig("")
	if err != nil {
		panic(err)
	}

	app, err := app.New(ctx, cfg, log)
	if err != nil {
		panic(err)
	}

	go func() {
		app.Start(ctx)
	}()
	log.Info("server is running")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	app.Stop(ctx)
	log.Info("server shutdown")
}
