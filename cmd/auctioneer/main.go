package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	conf "github.com/dog-sky/auctioneer/configs"
	server "github.com/dog-sky/auctioneer/internal/app/auctioneer"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	cfg, err := conf.NewConfig()
	if err != nil {
		log.Fatalf("Error loading conf: %v", err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	app, err := server.NewApp(ctx, cfg)
	if err != nil {
		log.Fatalf("Error setting up server: %v", err)
	}
	app.Setup()

	if err = app.MakeBlizzAuth(); err != nil {
		log.Fatal(err)
	}
	if err = app.GetRealmList(); err != nil {
		log.Fatal(err)
	}

	go func() {
		oscall := <-interrupt
		log.Printf("Syscall: %+v", oscall)
		cancel()
	}()

	app.Serve()
}
