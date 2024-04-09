package handler

import (
	"net/http"

	"github.com/carry0987/FileTree-API/internal/security"
	"github.com/carry0987/FileTree-API/internal/service"
	"github.com/carry0987/FileTree-API/internal/utils"
	"github.com/carry0987/FileTree-API/pkg/api"
	"github.com/gorilla/mux"
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

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.OutputMessage(w, utils.WebSocketResponse, http.StatusInternalServerError, "Failed to upgrade to WebSocket")
		return
	}
	defer conn.Close()

	// Get the signature and encrypted parameters from the route or query parameters
	vars := mux.Vars(r)
	encryptedPath := vars["encrypted"]
	if encryptedPath == "" {
		errorMsg := "Missing encrypted parameter"
		api.UnauthorizedError(errorMsg)
		utils.OutputMessage(w, utils.WebSocketResponse, http.StatusBadRequest, errorMsg)
		return
	}

	decryptedPath, err := security.Decrypt(encryptedPath)
	if err != nil {
		utils.OutputMessage(w, utils.WebSocketResponse, http.StatusInternalServerError, "Failed to decrypt the path")
		return
	}

	realPath, organize := utils.CheckOrganize(decryptedPath)
	fileTreeResult, err := service.GenerateFileTree(realPath, organize)
	if err != nil {
		utils.OutputMessage(w, utils.WebSocketResponse, http.StatusInternalServerError, "Error generating file tree")
		return
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	result, err := json.Marshal(fileTreeResult)
	if err != nil {
		utils.OutputMessage(w, utils.WebSocketResponse, http.StatusInternalServerError, "Error encoding file tree result to JSON")
		return
	}

	if err = conn.WriteMessage(websocket.TextMessage, result); err != nil {
		utils.OutputMessage(w, utils.WebSocketResponse, http.StatusInternalServerError, "Failed to send file tree result over WebSocket")
		return
	}
}
