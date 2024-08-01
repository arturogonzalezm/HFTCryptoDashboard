package handlers

import (
	"HFTCryptoDashboard/internal/websocket"
	"log"
	"os"
	"time"

	ws "github.com/gorilla/websocket"
)

func HandleInterrupt(interrupt chan os.Signal, done chan struct{}, client *websocket.WebSocketClient) {
	select {
	case <-interrupt:
		log.Println("Interrupt received, shutting down")
	case <-done:
	}

	err := client.Conn.WriteMessage(ws.CloseMessage, ws.FormatCloseMessage(ws.CloseNormalClosure, ""))
	if err != nil {
		log.Fatalf("Failed to send close message: %v", err)
	}
	select {
	case <-done:
	case <-time.After(time.Second):
	}
}
