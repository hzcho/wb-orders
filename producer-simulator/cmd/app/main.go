package main

import (
	"context"
	"os"
	"os/signal"
	"producer-simulator/internal/app"
	"producer-simulator/internal/config"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	log := logrus.New()

	godotenv.Load(".env.producer.example")

	cfg, err := config.InitConfig("")
	if err != nil {
		panic(err)
	}

	app := app.New(ctx, cfg, log)

	go func() {
		app.Start(ctx)
	}()
	log.Info("server is running")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	app.Stop()
	log.Info("server shutdown")
}
