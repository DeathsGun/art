package api

import (
	"fmt"
	"os"
	"testing"
)

func TestLogin(t *testing.T) {
	untis, err := NewUntisAPI("bk-ahaus")
	err = untis.Login("schm132074", os.Getenv("UNTIS_PW"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(untis.loginResponse.SessionId)
	err = untis.Logout()
	if err != nil {
		t.Fatal(err)
	}
}
