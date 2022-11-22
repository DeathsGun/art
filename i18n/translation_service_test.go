package i18n

import (
	"context"
	"testing"
)

var translationService = New()

func TestService_Translate(t *testing.T) {
	tests := []struct {
		language string
		key      string
		value    string
	}{
		{
			"en",
			"USERNAME",
			"Username",
		},
		{
			"de-ch",
			"USERNAME",
			"Benutzername",
		},
		{
			"fr-fr",
			"USERNAME",
			"USERNAME",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.language+"/"+test.key, func(t *testing.T) {
			t.Parallel()
			translationService.Translate(context.WithValue(context.Background(), LanguageCtxKey, test.language), test.key)
		})
	}
}
