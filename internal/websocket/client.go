package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"HFTCryptoDashboard/internal/config"
	"HFTCryptoDashboard/pkg/db"
	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	Conn   *websocket.Conn
	config *config.Config
	mu     sync.Mutex
}

var instance *WebSocketClient
var once sync.Once
var onceErr error

func GetWebSocketClient(cfg *config.Config) (*WebSocketClient, error) {
	once.Do(func() {
		conn, err := connectAndSubscribe("wss://stream.binance.com:9443/ws", cfg.Symbols)
		if err != nil {
			onceErr = err
			return
		}
		instance = &WebSocketClient{
			Conn:   conn,
			config: cfg,
		}
	})

	if onceErr != nil {
		return nil, onceErr
	}
	return instance, nil
}

func connectAndSubscribe(url string, symbols []string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	for _, symbol := range symbols {
		message := struct {
			Method string   `json:"method"`
			Params []string `json:"params"`
			ID     int      `json:"id"`
		}{
			Method: "SUBSCRIBE",
			Params: []string{symbol + "@ticker"},
			ID:     1,
		}
		err := conn.WriteJSON(message)
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

		var ticker db.TickerData
		err = json.Unmarshal(message, &ticker)
		if err != nil {
			log.Printf("Error unmarshalling message: %v, message: %s", err, message)
			continue
		}

		handleTickerData(ticker)
	}
}

func (client *WebSocketClient) CloseConnection() {
	client.mu.Lock()
	defer client.mu.Unlock()
	if err := client.Conn.Close(); err != nil {
		log.Printf("Failed to close WebSocket connection: %v", err)
	}
}

func handleTickerData(ticker db.TickerData) {
	if ticker.Symbol == "" || ticker.EventTime == 0 {
		log.Printf("Invalid ticker data received: %+v", ticker)
		return
	}

	err := db.InsertTickerData(ticker)
	if err != nil {
		log.Printf("Failed to insert ticker data: %v", err)
	} else {
		log.Printf("Successfully inserted ticker data: %+v", ticker)
	}
}
