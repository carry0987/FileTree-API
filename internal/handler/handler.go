package handler

import (
	"net/http"

	"github.com/carry0987/FileTree-API/internal/utils"
	"github.com/gorilla/websocket"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	utils.OutputMessage(w, utils.HTTPResponse, http.StatusOK, "FileTree API")
}

func UnifiedHandler(w http.ResponseWriter, r *http.Request) {
	if websocket.IsWebSocketUpgrade(r) {
		WebSocketHandler(w, r)
	} else {
		HTTPHandler(w, r)
	}
}
