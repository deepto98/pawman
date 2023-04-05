package cmd

import (
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

const (
	KeySize   = 32
	NonceSize = 24
)

type EncryptPayload struct {
	EncryptionKey *[KeySize]byte
	Nonce         *[NonceSize]byte
	Ciphertext    []byte
}

// Func to generate key
func GenerateKey() (*[KeySize]byte, error) {
	key := new([KeySize]byte)
	// Generate random no and copy to str
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		return nil, err
	}
	return key, nil
}

// Func to generate nonce
func GenerateNonce() (*[NonceSize]byte, error) {
	nonce := new([NonceSize]byte)
	// Generate random no and copy to str
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return nil, err
	}
	return nonce, nil
}

func Encrypt(message []byte) ([]byte, *[KeySize]byte, *[NonceSize]byte, error) {
	nonce, err := GenerateNonce()
	if err != nil {
		return nil, nil, nil, errors.New("Unable to generate nonce")
	}
	key, err := GenerateKey()
	if err != nil {
		return nil, nil, nil, errors.New("Unable to generate key")
	}

	// Generate Ciphertext
	ciphertext := make([]byte, len(nonce))
	copy(ciphertext, nonce[:])
	ciphertext = secretbox.Seal(ciphertext, message, nonce, key)

	return ciphertext, key, nonce, nil
}
