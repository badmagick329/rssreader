package auth

import (
	"net/http"
	"strings"
)

type NoApiKeyError struct{}

func (e NoApiKeyError) Error() string {
	return "No API key provided"
}

type InvalidApiKeyError struct{}

func (e InvalidApiKeyError) Error() string {
	return "Invalid API key format"
}

const API_KEY_HEADER = "Authorization"

// GetApiKey returns the API key from the request header
// Example:
// Authorization: ApiKey API_KEY
func GetApiKey(h http.Header) (string, error) {
	apiKeyHeader := h.Get(API_KEY_HEADER)
	if apiKeyHeader == "" {
		return "", NoApiKeyError{}
	}
	vals := strings.Split(apiKeyHeader, " ")
	if len(vals) != 2 {
		return "", InvalidApiKeyError{}
	}
	if vals[0] != "ApiKey" {
		return "", InvalidApiKeyError{}
	}
	return vals[1], nil
}
