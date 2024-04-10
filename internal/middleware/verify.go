package middleware

import (
	"fmt"
	"net/http"

	"github.com/carry0987/FileTree-API/internal/handler"
	"github.com/carry0987/FileTree-API/internal/security"
	"github.com/carry0987/FileTree-API/internal/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func SignatureVerificationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the signature and encrypted parameters from the route or query parameters
		vars := mux.Vars(r)
		signature := vars["signature"]
		encrypted := vars["encrypted"]

		// Reassemble the URL path for signature verification
		encPath := fmt.Sprintf("/enc/%s", encrypted)

		// Decode signature
		_, err := utils.Base64UrlDecode(signature)
		if err != nil {
			if websocket.IsWebSocketUpgrade(r) {
				handler.WebSocketMessage(w, r, handler.ErrInvalidSignatureFormat.Error())
			} else {
				utils.OutputMessage(w, utils.HTTPResponse, http.StatusBadRequest, handler.ErrInvalidSignatureFormat.Error())
			}
			return
		}

		// Verify the signature
		if !security.VerifySignature(signature, encPath) {
			if websocket.IsWebSocketUpgrade(r) {
				handler.WebSocketMessage(w, r, handler.ErrInvalidSignature.Error())
			} else {
				utils.OutputMessage(w, utils.HTTPResponse, http.StatusForbidden, handler.ErrInvalidSignature.Error())
			}
			return
		}

		// Continue processing the rest of the request
		next.ServeHTTP(w, r)
	})
}
