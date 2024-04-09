package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	message := "FileTree API"
	if websocket.IsWebSocketUpgrade(r) {
		WebSocketMessage(w, r, message)
	} else {
		HTTPMessage(w, r, message)
	}
}

func UnifiedHandler(w http.ResponseWriter, r *http.Request) {
	if websocket.IsWebSocketUpgrade(r) {
		WebSocketHandler(w, r)
	} else {
		HTTPHandler(w, r)
	}
}
