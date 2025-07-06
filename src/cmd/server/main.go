// main.go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spirozh/timr/internal"
)

func main() {
	a := internal.NewApp()
	go a.Run()

	// Set up signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	if err := a.Close(); err != nil {
		log.Fatal("Error closing server:", err)
	}
	log.Println("Server closed")
}
