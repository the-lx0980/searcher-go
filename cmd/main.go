package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourusername/wroxen-go/internal/bot"
	"github.com/yourusername/wroxen-go/internal/config"
	"github.com/yourusername/wroxen-go/internal/user"
)

func main() {
	// load env (config package uses godotenv automatically in its init)
	cfg := config.Get() // returns pointer to Config

	// Create user (TDLib) client
	userClient := user.NewUser(cfg)

	// Start user client (makes sure TDLib client initialized)
	_, err := userClient.Start()
	if err != nil {
		// non-fatal: Print warning; TDLib may still be usable if session already exists
		log.Printf("user start warning: %v\n", err)
	}

	// Create bot
	wroxen, err := bot.NewWroxen(cfg, userClient)
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	// Run bot in goroutine
	go wroxen.Start()

	// Wait for interrupt to stop gracefully
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down...")
	wroxen.Stop()
}
