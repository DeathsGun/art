package auth

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/deathsgun/art/untis"
	"io"
	"os"
)

var key = []byte(os.Getenv("CRYPT_KEY"))

func DecryptSession(rawSession string) (*untis.Session, error) {
	cipherData, err := base64.RawStdEncoding.DecodeString(rawSession)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(cipherData) < aes.BlockSize {
		return nil, errors.New("data block size is too short")
	}

	iv := cipherData[:aes.BlockSize]
	cipherData = cipherData[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherData, cipherData)

	session := &untis.Session{}
	err = json.Unmarshal(cipherData, session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func EncryptSession(session *untis.Session) (string, error) {
	data, err := json.Marshal(session)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	cipherData := make([]byte, aes.BlockSize+len(data))
	iv := cipherData[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherData[aes.BlockSize:], data)
	return base64.RawStdEncoding.EncodeToString(cipherData), err
}

func Session(ctx context.Context) *untis.Session {
	return ctx.Value("session").(*untis.Session)
}
