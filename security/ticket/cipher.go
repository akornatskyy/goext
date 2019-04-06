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
	b, _ := aes.NewCipher(c.key)
	gcm, _ := cipher.NewGCM(b)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, src, nil)
	return base64.RawURLEncoding.EncodeToString(ciphertext), nil
}

func (c *c) DecodeString(s string) ([]byte, error) {
	data, _ := base64.RawURLEncoding.DecodeString(s)
	b, _ := aes.NewCipher(c.key)
	gcm, _ := cipher.NewGCM(b)
	n := gcm.NonceSize()
	nonce, ciphertext := data[:n], data[n:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
