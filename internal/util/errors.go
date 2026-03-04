package util

import (
	"encoding/json"
	"fmt"
	"os"
)

// ErrorResponse is an LLM-friendly error format output as JSON to stderr.
type ErrorResponse struct {
	Error   bool   `json:"error"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

// ExitWithError outputs a JSON error to stderr and exits with code 1.
func ExitWithError(errType, message string) {
	resp := ErrorResponse{
		Error:   true,
		Type:    errType,
		Message: message,
	}
	data, _ := json.Marshal(resp)
	fmt.Fprintln(os.Stderr, string(data))
	os.Exit(1)
}

// FormatError returns a JSON error string without exiting.
func FormatError(errType, message string) string {
	resp := ErrorResponse{
		Error:   true,
		Type:    errType,
		Message: message,
	}
	data, _ := json.Marshal(resp)
	return string(data)
}
