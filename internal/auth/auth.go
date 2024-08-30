package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Extract APIKey from headers else error
// Example:
// Header Key: Authorization ApiKey {api key itself}

func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")

	if val == "" {
		return "", errors.New("no authentication headers")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed authorization headers")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("malformed first section of auth header")
	}

	return vals[1], nil
}
