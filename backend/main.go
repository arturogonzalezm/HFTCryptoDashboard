package main

import (
	"HFTCryptoDashboard/internal/config"
	"HFTCryptoDashboard/internal/handlers"
	"HFTCryptoDashboard/internal/websocket"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client, err := websocket.GetWebSocketClient(cfg)
	if err != nil {
		log.Fatalf("Failed to get WebSocket client: %v", err)
	}
	defer client.CloseConnection()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})
	go client.ReadMessages(done)

	handlers.HandleInterrupt(interrupt, done, client)
}
