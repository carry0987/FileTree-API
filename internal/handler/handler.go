package handler

import (
	"errors"
	"net/http"

	"github.com/carry0987/FileTree-API/internal/security"
	"github.com/carry0987/FileTree-API/internal/service"
	"github.com/carry0987/FileTree-API/internal/utils"
	"github.com/carry0987/FileTree-API/pkg/api"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Error declaration for specific Error handling
var (
	ErrMissingEncryptedParam = errors.New("missing encrypted parameter")
	ErrFailedToDecrypt       = errors.New("failed to decrypt the path")
	ErrErrorGeneratingFileTree = errors.New("error generating file tree")
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

// Processes the encrypted path, decrypts it, and generates the file tree.
func ProcessEncryptedPath(r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	// Get the signature and encrypted parameters from the route or query parameters
	encryptedPath := vars["encrypted"]
	if encryptedPath == "" {
		api.UnauthorizedError(ErrMissingEncryptedParam.Error())
		return nil, ErrMissingEncryptedParam
	}

	// Decrypt the path
	decryptedPath, err := security.Decrypt(encryptedPath)
	if err != nil {
		api.UnauthorizedError(ErrFailedToDecrypt.Error())
		return nil, ErrFailedToDecrypt
	}

	// Generate the file tree using the decrypted path
	realPath, organize := utils.CheckOrganize(decryptedPath)
	fileTreeResult, err := service.GenerateFileTree(realPath, organize)
	if err != nil {
		api.InternalServerError(ErrErrorGeneratingFileTree.Error())
		return nil, ErrErrorGeneratingFileTree
	}

	return fileTreeResult, nil
}

// Maps specific error types to HTTP status codes
func DetermineHTTPStatusCode(err error) int {
	switch err {
	case ErrMissingEncryptedParam:
		// Missing parameter
		return http.StatusUnauthorized
	case ErrFailedToDecrypt:
		// Failed to decrypt could imply a wrong input
		return http.StatusUnauthorized
	case ErrErrorGeneratingFileTree:
		// Error generating file tree implies internal server problems
		return http.StatusInternalServerError
	default:
		// For any other error, consider it as an internal server error
		return http.StatusInternalServerError
	}
}
