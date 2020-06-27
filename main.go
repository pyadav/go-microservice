package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/pyadav/microservice/handlers"
)

func main() {
	logger := log.New(os.Stdout, "product-api", log.LstdFlags)

	productsHandler := handlers.NewProducts(logger)

	serverMux := http.NewServeMux()
	serverMux.Handle("/", productsHandler)

	server := &http.Server{
		Addr:         ":9090",
		Handler:      serverMux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, os.Kill)

	signal := <-signalChan
	logger.Println("Recived terminate, graceful shutdown", signal)

	timeoutContext, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(timeoutContext)
}
