package cmd

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"io/ioutil"

	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/ssh"
)

func DecryptWithPrivateKey(privKeyLoc string, wrappedData []byte) ([]byte, error) {
	// 1.Read priv key
	priv, err := ioutil.ReadFile(privKeyLoc)
	if err != nil {
		return nil, errors.New("Unable to read priv key")
	}
	// Parse key
	privKey, err := ssh.ParseRawPrivateKey(priv)
	if err != nil {
		return nil, errors.New("Unable to parse priv key")
	}

	// 2. Get raw encrypted payload
	data, err := base64.StdEncoding.DecodeString(string(wrappedData))
	if err != nil {
		return nil, errors.New("Unable to decode payload")
	}

	// Parse OpenSSH key as RSA Private key
	parsedKey := privKey.(*rsa.PrivateKey)
	// Decrypt
	decryptedBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, parsedKey, data, nil)
	if err != nil {
		return nil, errors.New("Unable to decrypt RSA key")
	}
	return decryptedBytes, nil

}

func DecryptBox(gobText bytes.Buffer) ([]byte, error) {
	//Decode gob
	gobs := gob.NewDecoder(&gobText)
	var encPayload EncryptPayload
	err := gobs.Decode(&encPayload)
	if err != nil {
		return nil, err
	}
	//Get Nonce, Key and Ciphertext from Gob
	var nonce [24]byte
	copy(nonce[:], encPayload.Ciphertext[:24])

	//Decrypt with Secretbox
	out, ok := secretbox.Open(nil, encPayload.Ciphertext[24:], encPayload.Nonce, encPayload.EncryptionKey)
	if !ok {
		return nil, errors.New("unable to decrypt")
	}

	return out, nil

}
