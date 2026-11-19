package certmanager

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// The API answers HTTP 200 with header.isSuccessful=false on failure.
func TestFailedHeaderIsAnError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"header":{"isSuccessful":false,"resultCode":4010001,"resultMessage":"Invalid appKey."}}`))
	}))
	defer srv.Close()
	c := NewClient("example-app-key", "", "", nil, false)
	c.baseURL = srv.URL
	_, err := c.ListCertificates(context.Background())
	if err == nil || !strings.Contains(err.Error(), "4010001") || !strings.Contains(err.Error(), "Invalid appKey.") {
		t.Fatalf("want an error with the result code and message, got %v", err)
	}
}

// The v1.3 list answers dates as yyyy-mm-dd and the ID as a number (API
// guide response example); the whole list must decode.
func TestListCertificatesDecodesGuideExample(t *testing.T) {
	var query string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.Path + "?" + r.URL.RawQuery
		w.Write([]byte(`{"header":{"resultCode":0,"resultMessage":"success","isSuccessful":true},
"body":{"totalCount":1,"totalPage":1,"currentPage":1,"pageSize":10,"data":[{"certificateId":1,
"certificateName":"example.com","authority":"NHN","domains":["example.com","*.example.com"],
"signatureAlgorithm":"SHA256withRSA","fileCreationDate":"2025-03-02","expirationDate":"2026-03-25"}]}}`))
	}))
	defer srv.Close()
	c := NewClient("example-app-key", "", "", nil, false)
	c.baseURL = srv.URL
	out, err := c.ListCertificates(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if query != "/certmanager/v1.3/appkeys/example-app-key/certificates?all=true" {
		t.Errorf("request %s", query)
	}
	got := out.Body.Certificates
	if len(got) != 1 || got[0].CertificateID != 1 || got[0].CertificateName != "example.com" || got[0].Authority != "NHN" ||
		len(got[0].Domains) != 2 || got[0].ExpirationDate != "2026-03-25" || got[0].FileCreationDate != "2025-03-02" {
		t.Errorf("decoded %+v", got)
	}
}
