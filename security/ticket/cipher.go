// Package ticket provides ticket signing and encryption utilities for secure token generation.
package ticket

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

type c struct {
	key []byte
}

func NewCipher(key []byte) *c {
	return &c{
		key: key,
	}
}

func (c *c) EncodeToString(src []byte) (string, error) {
	b, err := aes.NewCipher(c.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(b)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, src, nil)
	return base64.RawURLEncoding.EncodeToString(ciphertext), nil
}

func (c *c) DecodeString(s string) ([]byte, error) {
	data, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	b, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(b)
	if err != nil {
		return nil, err
	}
	n := gcm.NonceSize()
	nonce, ciphertext := data[:n], data[n:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
