package haclient

import "fmt"

// APIError represents an error from the Home Assistant API.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("HA API error (status %d): %s", e.StatusCode, e.Message)
}

// IsNotFound returns true if the error is a 404.
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == 404
}
