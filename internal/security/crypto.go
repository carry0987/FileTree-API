package security

import (
	"FileTree-API/internal/utils"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

// key and salt should be retrieved from environment variables or other configuration sources here.
var (
	key  []byte
	salt []byte
)

func SetKeyAndSalt(k []byte, s []byte) {
	key = k
	salt = s
}

// Decrypt is used to decrypt the encrypted message passed in from the client
func Decrypt(encryptedMessage string) (string, error) {
	cipherText, err := utils.Base64UrlDecode(encryptedMessage)
	if err != nil {
		return "", err
	}

	if len(cipherText) < 12 { // GCM nonce size expected to be 12 bytes
		return "", errors.New("cipherText too short")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	nonce := cipherText[:12]
	cipherText = cipherText[12:]

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aead.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// Compare the HMAC created from the message and salt with the one provided by the client
func VerifySignatures(signature, encryptedPath string) bool {
	// Decode the signature to get the HMAC
	signatureDecoded, err := utils.Base64UrlDecode(signature)
	if err != nil {
		utils.OutputMessage(nil, utils.LogOutput, 0, "Invalid: %v", signature)
		return false
	}

	// Compute the HMAC for the encryptedPath with the key and salt
	mac := hmac.New(sha256.New, key)
	mac.Write(salt)
	mac.Write([]byte(encryptedPath))
	expectedMAC := mac.Sum(nil)

	// Compare the client's HMAC with the expected HMAC
	return hmac.Equal(signatureDecoded, expectedMAC)
}

func VerifySignature(signature, encryptedPath string) bool {
	decodedSignature, err := utils.Base64UrlDecode(signature)
	if err != nil {
		fmt.Printf("Invalid signature decode error: %v\n", err)
		return false
	}

	mac := hmac.New(sha256.New, key)
	mac.Write(salt)
	mac.Write([]byte(encryptedPath))
	expectedMAC := mac.Sum(nil)

	return hmac.Equal(decodedSignature, expectedMAC)
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

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
			utils.OutputMessage(w, utils.HTTPResponse, http.StatusBadRequest, "Invalid signature format")
			return
		}

		// Verify the signature
		if !VerifySignature(signature, encPath) {
			utils.OutputMessage(w, utils.HTTPResponse, http.StatusForbidden, "Invalid signature")
			return
		}

		// Continue processing the rest of the request
		next.ServeHTTP(w, r)
	})
}
