package main

import (
	"github.com/SakuraBurst/vigilant-octo-meme/internal/app"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/config"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log := logger.NewSlogLogger(cfg)
	application, err := app.NewApp(cfg, log)
	if err != nil {
		panic(err)
	}
	stop := make(chan os.Signal, 1)
	go func() {
		err := application.Run()
		if err != nil {
			log.Error(err.Error())
			stop <- syscall.SIGTERM
		}
	}()

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	application.Stop()
}
