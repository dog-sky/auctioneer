package main

import (
	server "auctioneer/app/auctioneer"
	"auctioneer/app/conf"
	"context"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := conf.NewConfig()
	if err != nil {
		log.Fatalf("Error loading conf: %v", err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	app, err := server.Setup(ctx, cfg)
	if err != nil {
		log.Fatalf("Error setting up server: %v", err)
	}

	go func() {
		oscall := <-interrupt
		log.Printf("Syscall: %+v", oscall)
		cancel()
	}()

	app.Serve()
}
