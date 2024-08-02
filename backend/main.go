package main

import (
	"HFTCryptoDashboard/internal/config"
	"HFTCryptoDashboard/internal/handlers"
	"HFTCryptoDashboard/internal/websocket"
	"HFTCryptoDashboard/pkg/db"
	"HFTCryptoDashboard/pkg/util"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		util.LogAndExit("Error loading .env file", 1)
	}

	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		util.LogAndExit("Failed to load config: %v", 1, err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		util.LogAndExit("Database configuration environment variables are missing", 1)
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Initialize the database
	db.InitDB(dbURL)

	client, err := websocket.GetWebSocketClient(cfg)
	if err != nil {
		util.LogAndExit("Failed to get WebSocket client: %v", 1, err)
	}
	if client == nil {
		util.LogAndExit("WebSocket client is nil", 1)
	}
	defer client.CloseConnection()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})
	go func() {
		client.ReadMessages(done)
	}()

	handlers.HandleInterrupt(interrupt, done, client)
}
