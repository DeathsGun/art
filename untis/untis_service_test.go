package untis

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var untis = NewService()
var username = os.Getenv("UNTIS_USERNAME")
var password = os.Getenv("UNTIS_PASSWORD")
var school = os.Getenv("UNTIS_SCHOOL")

func TestService_Login(t *testing.T) {
	schools, err := untis.SearchSchools(context.Background(), school)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotEmpty(t, schools) {
		return
	}
	session, err := untis.Login(context.Background(), &schools[0], username, password)
	if !assert.NoError(t, err) {
		return
	}
	t.Cleanup(func() {
		_ = untis.Logout(context.Background(), session)
	})
	if !assert.NotEmpty(t, session.SessionId) {
		return
	}
	if !assert.NotEmpty(t, session.Endpoint) {
		return
	}
	if !assert.NotEmpty(t, session.PersonType) {
		return
	}
	if !assert.NotEmpty(t, session.PersonId) {
		return
	}

	err = untis.ValidateLogin(context.Background(), session)
	if !assert.NoError(t, err) {
		return
	}
}
