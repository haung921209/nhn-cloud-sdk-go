package nhncloud

import (
	"net/http"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
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
		return c.HTTPClient
	}
	return http.DefaultClient
}

func (c *Config) userAgent() string {
	if c.UserAgent != "" {
		return c.UserAgent
	}
	return "nhn-cloud-sdk-go/0.1.0"
}
