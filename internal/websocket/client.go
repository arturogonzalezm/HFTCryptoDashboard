package websocket

import (
	"fmt"
	"log"
	"sync"

	"HFTCryptoDashboard/internal/config"
	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	Conn   *websocket.Conn
	config *config.Config
	mu     sync.Mutex
}

var instance *WebSocketClient
var once sync.Once

func GetWebSocketClient(cfg *config.Config) (*WebSocketClient, error) {
	var err error
	once.Do(func() {
		conn, err := connectAndSubscribe("wss://stream.binance.com:9443/ws", cfg.Symbols)
		if err != nil {
			log.Fatalf("Failed to connect and subscribe: %v", err)
		}
		instance = &WebSocketClient{
			Conn:   conn,
			config: cfg,
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

func (client *WebSocketClient) ReadMessages(done chan struct{}) {
	defer close(done)
	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}
		log.Printf("Received message: %s", message)
	}
}

func (client *WebSocketClient) CloseConnection() {
	client.mu.Lock()
	defer client.mu.Unlock()
	if err := client.Conn.Close(); err != nil {
		log.Printf("Failed to close WebSocket connection: %v", err)
	}
}
