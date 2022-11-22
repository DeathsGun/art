package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/deathsgun/art/crypt"
	"github.com/deathsgun/art/di"
	"github.com/deathsgun/art/untis"
)

func DecryptSession(rawSession string) (*untis.Session, error) {
	data, err := base64.StdEncoding.DecodeString(rawSession)
	if err != nil {
		return nil, err
	}
	data, err = di.Instance[crypt.ICryptService]("crypt").Decrypt(data)
	if err != nil {
		return nil, err
	}

	session := &untis.Session{}
	err = json.Unmarshal(data, session)
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
	data, err = di.Instance[crypt.ICryptService]("crypt").Encrypt(data)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

func Session(ctx context.Context) *untis.Session {
	return ctx.Value("session").(*untis.Session)
}
