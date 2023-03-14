package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/likekimo/goproxyspider/container"
	"gitlab.com/likekimo/goproxyspider/models"
	"gitlab.com/likekimo/goproxyspider/pkg/logger"
	"gitlab.com/likekimo/goproxyspider/spiders"
	"gitlab.com/likekimo/goproxyspider/storage"
	"go.uber.org/zap"
)

func main() {
	var (
		toSave   = flag.String("to.save", "badgerdb", "To save proxy default proxy.txt")
		logLevel = flag.String("verbose", "debug", "debug,info,warn,error,dpanic,panic,fatal")
		cmd      = flag.String("run", "init", "init,check")
	)
	flag.Parse()

	logger := logger.NewLogger(*logLevel)
	logger.Info("hello")
	defer logger.Info("goodbye")

	errc := make(chan error)

	ctx := context.Background()
	proxys := make(chan models.Proxy)

	// conf storage
	var writeStorage storage.Storage
	if *toSave == "" {
		writeStorage = storage.FileStorage(logger, "")
	} else if *toSave == "badgerdb" {
		badgerStorage := storage.BadgerStorage(logger, "")
		err := badgerStorage.CheckStorage()
		if err != nil {
			logger.Error(err.Error())
			errc <- fmt.Errorf("%s", err)
		}
		//defer badgerStorage.
		writeStorage = badgerStorage
	}

	if *cmd == "init" {
		spdSlice := container.New(logger, proxys)

		for _, v := range spdSlice {
			go v.Run(ctx)
		}

		workers := len(spdSlice)
		go func() {
			for {
				msg := <-proxys
				if msg.Source == "stop" { // все останавливаем
					workers = workers - 1
					if workers == 0 {
						// выводим стату
						writeStorage.Stat()
						// end
						errc <- spiders.ProxySpiderError("End")
					}
					logger.Info("workers", zap.Int("workers", workers))
				}
				writeStorage.WriteStorage(msg.IP)
			}
		}()
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	logger.Info("exit", zap.Error(<-errc))
}
