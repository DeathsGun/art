package api

import (
	"github.com/deathsgun/art/rpc"
	"github.com/google/uuid"
)

type Untis struct {
	client        *rpc.Client
	loginResponse *LoginResponse
}

func (u *Untis) Login(username string, password string) error {
	err := u.client.Call("authenticate", LoginRequest{Username: username, Password: password, Client: uuid.NewString()}, u.loginResponse)
	return err
}

func (u *Untis) Logout() error {
	return u.client.Call("logout", nil, nil)
}

func NewUntisAPI(schoolName string) (*Untis, error) {
	client, err := rpc.Dial("https://asopo.webuntis.com/WebUntis/jsonrpc.do?school=" + schoolName)
	if err != nil {
		return nil, err
	}
	return &Untis{client: client}, nil
}

type Person int

const (
	Teacher  Person = 2
	Students Person = 5
)

type LoginResponse struct {
	SessionId  string `json:"sessionId"`
	PersonType Person `json:"personType"`
	PersonId   int    `json:"personId"`
}

type LoginRequest struct {
	Username string `json:"user"`
	Password string `json:"password"`
	Client   string `json:"client"`
}
