package utils

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type MessageOutputMode int

const (
	// HTTPResponse means the message will be output as HTTP response
	HTTPResponse MessageOutputMode = iota
	// WebSocketResponse means the message will be output as WebSocket response
	WebSocketResponse
	// LogOutput means the message will be logged
	LogOutput
	// FatalOutput means the message will be logged and the program will be terminated
	FatalOutput
)

// LoadEnv loads the environment variables from .env file
func LoadEnv() {
	// If .env file does not exist, do nothing
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		return
	}
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %s\n", err)
	}
}

// Check if the required environment variables exist
func CheckEnvVariables(variables []string) {
	for _, varName := range variables {
		if value := os.Getenv(varName); value == "" {
			OutputMessage(nil, FatalOutput, 0, "Environment variable %s is not set or empty.", varName)
		}
	}
}

// OutputMessage provides message output, output to HTTP or log according to mode
func OutputMessage(w http.ResponseWriter, mode MessageOutputMode, statusCode int, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)

	switch mode {
	case HTTPResponse:
		if w == nil {
			log.Printf("http.ResponseWriter is nil\n")
			return
		}
		http.Error(w, message, statusCode)
	case WebSocketResponse:
		if w == nil {
			log.Printf("http.ResponseWriter is nil\n")
			return
		}
		w.Write([]byte(message))
	case LogOutput:
		log.Print(message)
	case FatalOutput:
		log.Fatalf(message)
	}
}

func Base64UrlEncode(data string) string {
	data = strings.ReplaceAll(data, "-", "+")
	data = strings.ReplaceAll(data, "_", "/")

	return base64.RawURLEncoding.EncodeToString([]byte(data))
}

// Base64UrlDecode encodes a byte slice to a base64 string
func Base64UrlDecode(data string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(data)
}

// Get the parameters separate by '::'
func CheckOrganize(data string) (string, bool) {
	s := strings.Split(data, "::")
	if len(s) != 2 {
		return s[0], false
	}
	if s[1] == "org" {
		return s[0], true
	}

	return s[0], false
}
