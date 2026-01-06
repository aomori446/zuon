package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

func Encrypt(password string, plaintext []byte) ([]byte, error) {
	salt := make([]byte, 8)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	block, err := newCipherBlock(password, salt)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	prefix := append(salt, nonce...)
	return gcm.Seal(prefix, nonce, plaintext, nil), nil
}

func Decrypt(password string, fullData []byte) ([]byte, error) {
	if len(fullData) < 8 {
		return nil, errors.New("data too short")
	}

	salt := fullData[:8]
	ciphertext := fullData[8:]

	block, err := newCipherBlock(password, salt)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce := ciphertext[:nonceSize]
	actualCipher := ciphertext[nonceSize:]

	return gcm.Open(nil, nonce, actualCipher, nil)
}

func newCipherBlock(password string, salt []byte) (cipher.Block, error) {
	key, err := pbkdf2.Key(sha256.New, password, salt, 4096, 32)
	if err != nil {
		return nil, err
	}
	return aes.NewCipher(key)
}
