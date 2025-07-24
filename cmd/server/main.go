package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"bobshop/internal/platform/config"
	"bobshop/internal/platform/response"
)

func main() {
	envPath := flag.String("env", ".env", "env file path")
	configPath := flag.String("config", "./configs/config.dev.yaml", "config file path")
	flag.Parse()
	cfg, err := config.LoadConfig(*envPath, *configPath)
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Configure toggles exposing backend error details (enable only in development).
	response.Configure(config.IsDevelopment())

	app, cleanup, err := buildApp(cfg)
	if err != nil {
		log.Fatalf("could not initialize server: %v", err)
	}
	defer cleanup()

	// Graceful shutdown (SIGINT, SIGTERM)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down gracefully...")
		cleanup()
		os.Exit(0)
	}()

	engine := app.Engine
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("Server is running on %s\n", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
