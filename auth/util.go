package auth

import (
	"context"
	"encoding/json"
	"github.com/deathsgun/art/crypt"
	"github.com/deathsgun/art/di"
	"github.com/deathsgun/art/untis"
)

func DecryptSession(rawSession string) (*untis.Session, error) {
	data, err := di.Instance[crypt.ICryptService]("crypt").DecryptString(rawSession)
	if err != nil {
		return nil, err
	}

	session := &untis.Session{}
	err = json.Unmarshal([]byte(data), session)
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
	return di.Instance[crypt.ICryptService]("crypt").EncryptString(string(data))
}

func Session(ctx context.Context) *untis.Session {
	return ctx.Value("session").(*untis.Session)
}
