package apiheader

import (
	"errors"
	"testing"
)

func TestCheck(t *testing.T) {
	for _, ok := range []string{``, `not json`, `{}`, `{"header":{"isSuccessful":true,"resultCode":0}}`, `{"header":{}}`} {
		if err := Check([]byte(ok)); err != nil {
			t.Errorf("Check(%q) = %v, want nil", ok, err)
		}
	}
	err := Check([]byte(`{"header":{"isSuccessful":false,"resultCode":51000,"resultMessage":"invalid app key"}}`))
	var e *Error
	if !errors.As(err, &e) || e.Code != "51000" || e.Message != "invalid app key" {
		t.Fatalf("Check(failed header) = %v, want code 51000 and the message", err)
	}
}
