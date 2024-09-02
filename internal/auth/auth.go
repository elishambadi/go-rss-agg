package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Extracts APIKey from headers of HTTP request
// Example input:
// Authorization: ApiKey {api_key}
func GetAPIKey(headers http.Header) (string, error) {
	auth_header := headers.Get("Authorization")

	if auth_header == "" {
		return "", errors.New("no authorization headers found")
	}

	values := strings.Split(auth_header, " ")
	if len(values) != 2 {
		return "", errors.New("authorization header malformed")
	}

	if values[0] != "ApiKey" {
		return "", errors.New("First path of Auth header malformed")
	}

	return values[1], nil
}
