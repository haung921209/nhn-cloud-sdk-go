// Package ptr provides helper functions for creating pointers to primitive types.
// This is useful for optional fields in API requests.
package ptr

// Int returns a pointer to the given int value.
func Int(v int) *int {
	return &v
}

// Bool returns a pointer to the given bool value.
func Bool(v bool) *bool {
	return &v
}

// String returns a pointer to the given string value.
func String(v string) *string {
	return &v
}

// Float64 returns a pointer to the given float64 value.
func Float64(v float64) *float64 {
	return &v
}
