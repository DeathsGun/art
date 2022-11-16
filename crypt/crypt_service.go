package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

var key = []byte(os.Getenv("CRYPT_KEY"))

type ICryptService interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

type service struct {
}

func (c *service) EncryptString(s string) (string, error) {
	cipherData, err := c.Encrypt([]byte(base64.RawStdEncoding.EncodeToString([]byte(s))))
	if err != nil {
		return "", err
	}
	return string(cipherData), nil
}

func (c *service) DecryptString(s string) (string, error) {
	data, err := c.Decrypt([]byte(s))
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(data), nil
}

func (c *service) Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cipherData := make([]byte, aes.BlockSize+len(data))
	iv := cipherData[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherData[aes.BlockSize:], data)
	return cipherData, nil
}

func (c *service) Decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(data) < aes.BlockSize {
		return nil, errors.New("data block size is too short")
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)
	return data, nil
}

func New() ICryptService {
	return &service{}
}
