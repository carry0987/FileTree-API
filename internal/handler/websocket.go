package handler

import (
	"encoding/base64"
	"net/http"

	"github.com/carry0987/FileTree-API/internal/utils"
	"github.com/carry0987/FileTree-API/pkg/api"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

type wrapChunkData struct {
	Index       int    `json:"index"`
	TotalChunks int    `json:"totalChunks"`
	Progress    int    `json:"progress"`
	Data        string `json:"data"`
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wrapChunks(chunk []byte, index, totalChunks int) ([]byte, error) {
	data := wrapChunkData{
		Index:       index,
		TotalChunks: totalChunks,
		Progress:    (index + 1) * 100 / totalChunks,
		Data:        base64.StdEncoding.EncodeToString(chunk),
	}

	return json.Marshal(data)
}

func sendInChunks(conn *websocket.Conn, data []byte, chunkSize int) error {
	totalChunks := len(data) / chunkSize
	if len(data)%chunkSize != 0 {
		totalChunks++
	}

	for i := 0; i < totalChunks; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunk := data[start:end]

		message, err := wrapChunks(chunk, i, totalChunks)
		if err != nil {
			return err
		}

		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return err
		}
	}

	return nil
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

	// Get the file tree result
	fileTreeResult, err := ProcessEncryptedPath(r)

	if err != nil {
		errorMsg := err.Error()
		apiErrorResponse := api.NewErrorResponse(errorMsg)
		errJSON, _ := json.Marshal(apiErrorResponse)
		conn.WriteMessage(websocket.TextMessage, errJSON)
		return
	}

	result, err := json.Marshal(fileTreeResult)
	if err != nil {
		utils.OutputMessage(conn, utils.WebSocketResponse, http.StatusInternalServerError, "Error encoding file tree result to JSON")
		return
	}

	chunkSize := 10240
	if err = sendInChunks(conn, result, chunkSize); err != nil {
		utils.OutputMessage(conn, utils.WebSocketResponse, http.StatusInternalServerError, "Failed to send file tree result over WebSocket in chunks")
		return
	}

	if err = conn.WriteMessage(websocket.TextMessage, result); err != nil {
		utils.OutputMessage(conn, utils.WebSocketResponse, http.StatusInternalServerError, "Failed to send file tree result over WebSocket")
		return
	}
}
