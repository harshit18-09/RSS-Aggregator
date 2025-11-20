package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extracts the API key from the HTTP headers
// Authorization: APIKey
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authorization info provided")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed authorization header")
	}
	if vals[0] != "APIKey" {
		return "", errors.New("malformed authorization header first part")
	}
	return vals[1], nil
}
