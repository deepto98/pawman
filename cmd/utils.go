package cmd

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"io"
	"io/ioutil"

	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/ssh"
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

func EncryptWithPublicKey(pubKeyLoc string, ciphertext bytes.Buffer) ([]byte, error) {
	// 1.Read and parse SSH pub key
	pub, err := ioutil.ReadFile(pubKeyLoc)
	if err != nil {
		return nil, errors.New("Unable to read pub key")
	}
	parsed, _, _, _, err := ssh.ParseAuthorizedKey(pub)
	if err != nil {
		return nil, errors.New("Unable to parse pub key")
	}

	// 2. Parse pub key into RSA format
	parsedCryptoKey := parsed.(ssh.CryptoPublicKey)
	pubCrypto := parsedCryptoKey.CryptoPublicKey()
	rsaPub := pubCrypto.(*rsa.PublicKey)

	// 3. Encrypt the binary
	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPub, ciphertext.Bytes(), nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return encryptedBytes, nil
	// 4. Base64 encode and store in db
}
