package main

import (
	"context"
	"fmt"
	// "os"
	"os/signal"
	"syscall"
	"time"

	"github.com/i474232898/chatserver/configs"
	"github.com/i474232898/chatserver/internal/app/repositories"
	"github.com/i474232898/chatserver/internal/app/server"
	"github.com/i474232898/chatserver/internal/app"
)

const (
	_shutdownPeriod      = 15 * time.Second
	_shutdownHardPeriod  = 3 * time.Second
	_readinessDrainDelay = 5 * time.Second
)

func main() {
	cfg := configs.New()

	_, dbErr := repositories.GetPool(cfg)
	if dbErr != nil {
		panic("Can't connect to db")
	}

	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srvr := server.NewServer()

	go func() {
		srvr.Start(ongoingCtx, cfg.Port)
	}()

	<-rootCtx.Done()
	stop()
	app.IsShuttingDown.Store(true)
	fmt.Println("Received termination signal, shutting down.")

	time.Sleep(_readinessDrainDelay)
	fmt.Println("Readiness check propagated, now waiting for ongoing requests to finish.")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), _shutdownPeriod)
	defer cancel()
	err := srvr.Shutdown(shutdownCtx)

	stopOngoingGracefully()
	
	if err != nil {
		fmt.Println("Failed to wait for ongoing requests to finish, waiting for forced cancellation.")
		time.Sleep(_shutdownHardPeriod)
	}

	fmt.Println("Server shut down gracefully.")
}
