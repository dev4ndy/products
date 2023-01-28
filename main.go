package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dev4ndy/products/handlers"
)

func main() {
	log := log.New(os.Stdout, "product-api", log.LstdFlags)

	log.Println("Starting Server on port 8080")
	productHandler := handlers.NewProducts(log)

	sm := http.NewServeMux()

	sm.Handle("/", productHandler)

	serve := http.Server{
		Addr:         ":8080",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	go func() {
		err := serve.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan

	log.Println("Recived terminatem, graceful shutdown", sig)

	timeoutContext, _ := context.WithTimeout(context.Background(), 30*time.Second)
	serve.Shutdown(timeoutContext)
}
