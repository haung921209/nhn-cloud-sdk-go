package nhncloud

import (
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/capture"
)

type Config struct {
	Region      string
	Credentials credentials.Credentials

	IdentityCredentials credentials.IdentityCredentials

	AppKeys map[string]string

	HTTPClient *http.Client
	Debug      bool
	UserAgent  string
}

func (c *Config) validate() error {
	if c.Region == "" {
		return ErrRegionRequired
	}
	if c.Credentials == nil {
		return ErrCredentialsRequired
	}
	return nil
}

func (c *Config) httpClient() *http.Client {
	if c.HTTPClient != nil {
		// Caller-supplied client: wrap its transport with capture so that
		// the opt-in NHN_SDK_CAPTURE_DIR mirroring still applies. Wrapping
		// is idempotent.
		return capture.WrapClient(c.HTTPClient)
	}
	// Use a fresh client (not http.DefaultClient) so we don't mutate the
	// global default transport for unrelated programs.
	return &http.Client{Transport: capture.NewTransport(http.DefaultTransport)}
}

// UserAgentString returns the user agent string for HTTP requests.
func (c *Config) UserAgentString() string {
	if c.UserAgent != "" {
		return c.UserAgent
	}
	return "nhn-cloud-sdk-go/0.1.0"
}
