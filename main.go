package main

import (
	"fmt"
	"os"
	"os/signal"

	"airport-transfer-app/internal/amadeus"
)

type app struct {
	amadeusClient *amadeus.Client
}

func main() {
	// Create a channel to handle the interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Start the application
	app := &app{
		amadeusClient: amadeus.New(),
	}
	startServer(app)

	// Wait for the interrupt signal before exiting
	<-interrupt
	fmt.Println("Exiting...")
}
