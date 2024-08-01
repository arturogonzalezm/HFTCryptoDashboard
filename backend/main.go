package main

import (
	"HFTCryptoDashboard/internal/config"
	"HFTCryptoDashboard/internal/handlers"
	"HFTCryptoDashboard/internal/websocket"
	"HFTCryptoDashboard/pkg/util" // Corrected import path
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		util.LogAndExit("Failed to load config: %v", 1, err)
	}

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
