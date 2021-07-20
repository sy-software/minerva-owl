package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
)

// AES256Encrypt encrypts a text using the given key bytes passed as hex string
func AES256Encrypt(key string, text string) (string, error) {
	keyBytes, _ := hex.DecodeString(key)
	textBytes := []byte(text)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// Never use more than 2^32 random nonces with a
	// given key because of the risk of a repeat.
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, textBytes, nil)
	final := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(final), nil
}

// AES256Decrypt decrypts a base64 text using the given key bytes passed as hex string
func AES256Decrypt(key string, text string) (string, error) {
	keyBytes, _ := hex.DecodeString(key)
	raw, err := base64.StdEncoding.DecodeString(text)

	if err != nil {
		return "", err
	}

	ciphertext := raw[12:]
	nonce := raw[0:12]

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
