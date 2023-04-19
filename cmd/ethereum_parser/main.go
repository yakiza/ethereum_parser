package main

import (
	"context"
	"ethereum_parser"
	"fmt"
	"github.com/caarlos0/env/v6"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Config struct {
	Address      string        `env:"HTTP_ADDRESS" envDefault:":8080"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"0"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"0"`
	IdleTimeout  time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"0"`
}

func main() {
	// Parse configuration coming from environment variables.
	var config Config
	if err := env.Parse(&config); err != nil {
		log.Fatal(err.Error())
	}

	// Setup interrupt handler. This optional step configures the process so
	// that SIGINT and SIGTERM signals cause the services to stop gracefully.
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)

	wg := new(sync.WaitGroup)
	defer wg.Wait()

	newSub := make(chan bool)

	// Initializing services
	var ethConfig ethereum_parser.EthereumClientConfig
	if err := env.Parse(&ethConfig); err != nil {
		log.Fatal(err.Error())
	}

	ethereumClient := ethereum_parser.NewEthereumClient(ethConfig)
	repo := ethereum_parser.NewMemStorage()
	service := ethereum_parser.NewService(&repo, ethereumClient, newSub)
	h := ethereum_parser.NewHTTPHandlers(&service)

	// Wiring up API
	mux := ethereum_parser.CreateAPIMux(h)

	// Starting HTTP server
	srv := http.Server{
		Handler:      mux,
		IdleTimeout:  config.IdleTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		Addr:         config.Address,
	}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	wg.Add(1)

	fmt.Printf("Running on %v\n", config.Address)

	// Start parsing service
	parserService := ethereum_parser.NewParserService(&repo, ethereumClient)
	go func() {
		err := parserService.Parse(context.Background(), newSub, wg)
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Add(1)

	// Graceful and eager terminations
	switch s := <-signalCh; s {
	case syscall.SIGTERM:
		log.Println("Terminating gracefully.")
		wg.Done()
		wg.Done()
		if err := srv.Shutdown(context.Background()); err != http.ErrServerClosed {
			log.Println("Failed to shutdown server:", err)
		}
	case syscall.SIGINT:
		log.Println("Terminating eagerly.")
		os.Exit(-int(syscall.SIGINT))
	}

	wg.Wait()
}
