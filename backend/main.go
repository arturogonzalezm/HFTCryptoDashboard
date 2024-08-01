package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

type Config struct {
	Symbols []string `yaml:"symbols"`
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	config, err := loadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to Binance WebSocket
	url := "wss://stream.binance.com:9443/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Failed to close WebSocket connection: %v", err)
		}
	}()

	// Subscribe to ticker streams
	for _, symbol := range config.Symbols {
		message := fmt.Sprintf("{\"method\": \"SUBSCRIBE\", \"params\": [\"%s@ticker\"], \"id\": 1}", symbol)
		err = conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Fatalf("Failed to send subscribe message: %v", err)
		}
	}

	// Set up a channel to handle interrupts for clean shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				return
			}
			log.Printf("Received message: %s", message)
		}
	}()

	select {
	case <-interrupt:
		log.Println("Interrupt received, shutting down")
	case <-done:
	}

	// Cleanly close the connection by sending a close message and then waiting for the server to close the connection
	err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Fatalf("Failed to send close message: %v", err)
	}
	select {
	case <-done:
	case <-time.After(time.Second):
	}
}
