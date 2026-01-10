package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Extracts the API key fro the http request headers
// Example: Authorization: ApiKey <apiKey>
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no authorization header provided")
	}

	authParts := strings.Split(authHeader, " ")
	if len(authParts) != 2 {
		return "", errors.New("malformed authorization header")
	}

	if authParts[0] != "ApiKey" {
		return "", errors.New("malformed authorization header, expected ApiKey")
	}

	return authParts[1], nil
}
