package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"os/signal"
	"sync"
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

type WebSocketClient struct {
	conn   *websocket.Conn
	config *Config
	mu     sync.Mutex
}

var instance *WebSocketClient
var once sync.Once

func GetWebSocketClient(config *Config) (*WebSocketClient, error) {
	var err error
	once.Do(func() {
		conn, err := connectAndSubscribe("wss://stream.binance.com:9443/ws", config.Symbols)
		if err != nil {
			log.Fatalf("Failed to connect and subscribe: %v", err)
		}
		instance = &WebSocketClient{
			conn:   conn,
			config: config,
		}
	})
	return instance, err
}

func connectAndSubscribe(url string, symbols []string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	for _, symbol := range symbols {
		message := fmt.Sprintf("{\"method\": \"SUBSCRIBE\", \"params\": [\"%s@ticker\"], \"id\": 1}", symbol)
		err = conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			return nil, err
		}
	}
	return conn, nil
}

func (client *WebSocketClient) readMessages(done chan struct{}) {
	defer close(done)
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}
		log.Printf("Received message: %s", message)
	}
}

func (client *WebSocketClient) closeConnection() {
	client.mu.Lock()
	defer client.mu.Unlock()
	if err := client.conn.Close(); err != nil {
		log.Printf("Failed to close WebSocket connection: %v", err)
	}
}

func main() {
	config, err := loadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client, err := GetWebSocketClient(config)
	if err != nil {
		log.Fatalf("Failed to get WebSocket client: %v", err)
	}
	defer client.closeConnection()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})
	go client.readMessages(done)

	select {
	case <-interrupt:
		log.Println("Interrupt received, shutting down")
	case <-done:
	}

	err = client.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Fatalf("Failed to send close message: %v", err)
	}
	select {
	case <-done:
	case <-time.After(time.Second):
	}
}
