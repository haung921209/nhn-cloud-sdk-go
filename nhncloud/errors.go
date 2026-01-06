package nhncloud

import "errors"

var (
	ErrRegionRequired      = errors.New("nhncloud: region is required")
	ErrCredentialsRequired = errors.New("nhncloud: credentials is required")
	ErrTenantIDRequired    = errors.New("nhncloud: tenant ID is required for this service")
	ErrAppKeyRequired      = errors.New("nhncloud: app key is required for this service")
)

type APIError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return "nhncloud: " + e.Code + ": " + e.Message
	}
	return "nhncloud: " + e.Message
}
