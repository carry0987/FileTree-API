package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"

	"FileTree-API/internal/utils"
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
func VerifySignature(signature, encryptedPath string) bool {
	// Decode the signature to get the HMAC
	decodedSignature, err := utils.Base64UrlDecode(signature)
	if err != nil {
		utils.OutputMessage(nil, utils.LogOutput, 0, "Invalid signature: %v\n", signature)
		return false
	}

	// Compute the HMAC for the encryptedPath with the key and salt
	mac := hmac.New(sha256.New, key)
	mac.Write(salt)
	mac.Write([]byte(encryptedPath))
	expectedMAC := mac.Sum(nil)

	// Compare the client's HMAC with the expected HMAC
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
