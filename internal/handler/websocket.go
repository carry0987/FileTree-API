package handler

import (
	"net/http"

	"github.com/carry0987/FileTree-API/internal/utils"
	"github.com/carry0987/FileTree-API/pkg/api"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func UpgradeToWebSocket(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        utils.OutputMessage(w, utils.WebSocketResponse, http.StatusInternalServerError, "Failed to upgrade to WebSocket")
        return nil, err
    }

    return conn, nil
}

func WebSocketMessage(w http.ResponseWriter, r *http.Request, message string) {
	conn, err := UpgradeToWebSocket(w, r)
	if err != nil {
		return
	}
	defer conn.Close()
	utils.OutputMessage(conn, utils.WebSocketResponse, http.StatusOK, message)
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := UpgradeToWebSocket(w, r)
	if err != nil {
		return
	}
	defer conn.Close()

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	fileTreeResult, err := ProcessEncryptedPath(r)

	if err != nil {
		errorMsg := err.Error()
		apiErrorResponse := api.NewErrorResponse(errorMsg)
		w.WriteHeader(DetermineHTTPStatusCode(err))
		conn.WriteJSON(apiErrorResponse)
		return
	}

	result, err := json.Marshal(fileTreeResult)
	if err != nil {
		utils.OutputMessage(conn, utils.WebSocketResponse, http.StatusInternalServerError, "Error encoding file tree result to JSON")
		return
	}

	if err = conn.WriteMessage(websocket.TextMessage, result); err != nil {
		utils.OutputMessage(conn, utils.WebSocketResponse, http.StatusInternalServerError, "Failed to send file tree result over WebSocket")
		return
	}
}
