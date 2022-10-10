package untis

import (
	"fmt"
	"github.com/deathsgun/art/rpc"
)

type Session struct {
	Endpoint   string `json:"endpoint"`
	SessionId  string `json:"sessionId"`
	School     string `json:"school"`
	Token      string `json:"token"`
	PersonType int    `json:"personType"`
	PersonId   int    `json:"personId"`
}

func (s *Session) NewClient() *rpc.Client {
	return rpc.New(s.Endpoint+"/jsonrpc.do?school="+s.School, s.SessionId)
}

func (s *Session) Id() string {
	return fmt.Sprintf("%s%d", s.School, s.PersonId)
}
