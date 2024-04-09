package main

import (
	"encoding/hex"
	"net/http"
	"os"

	"github.com/carry0987/FileTree-API/internal/handler"
	"github.com/carry0987/FileTree-API/internal/middleware"
	"github.com/carry0987/FileTree-API/internal/security"
	"github.com/carry0987/FileTree-API/internal/utils"

	"github.com/gorilla/mux"
)

var version = "1.1.5"

func main() {
	// Load the environment variables
	utils.LoadEnv()
	// Check if the required environment variables exist
	utils.CheckEnvVariables([]string{"FILETREE_SECRET_KEY", "FILETREE_SECRET_SALT"})

	// Decode the key and salt
	key, err := hex.DecodeString(os.Getenv("FILETREE_SECRET_KEY"))
	if err != nil {
		utils.OutputMessage(nil, utils.FatalOutput, 0, "Failed to decode FILETREE_SECRET_KEY")
	}
	salt, err := hex.DecodeString(os.Getenv("FILETREE_SECRET_SALT"))
	if err != nil {
		utils.OutputMessage(nil, utils.FatalOutput, 0, "Failed to decode FILETREE_SECRET_SALT")
	}

	// Pass the key and salt to the security package
	security.SetKeyAndSalt(key, salt)

	// Create a new Gorilla Mux HTTP router
	r := mux.NewRouter()

	// Add the gzip middleware to the router
	r.Use(middleware.GzipMiddleware)

	// Default handler for the root path
	r.Handle("/", http.HandlerFunc(handler.DefaultHandler))

	// Handler for WebSocket connections
	r.Handle("/ws", http.HandlerFunc(handler.WebSocketHandler))

	// Add the signature verification middleware to our file tree handler function
	r.Handle("/{signature}/enc/{encrypted}", security.SignatureVerificationMiddleware(http.HandlerFunc(handler.FileTreeHandler)))

	// If FILETREE_PORT is set, use that as the port, otherwise use 8080
	port := os.Getenv("FILETREE_PORT")
	if port == "" {
		port = "8080"
	}

	// Configure the server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start the server
	utils.OutputMessage(nil, utils.LogOutput, 0, "Listening on http://localhost%s\n", server.Addr)
	// Show the version
	utils.OutputMessage(nil, utils.LogOutput, 0, "Version: %s\n", version)
	// If the server fails to start, log the error
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		utils.OutputMessage(nil, utils.FatalOutput, 0, "ListenAndServe error: %v", err)
	}
}
