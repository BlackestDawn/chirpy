package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	auth := strings.TrimSpace(headers.Get("Authorization"))
	if auth == "" {
		return "", fmt.Errorf("no authorization header found")
	}

	tokenParts := strings.Split(auth, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		return "", fmt.Errorf("malformed bearer token: %s", auth)
	}

	return tokenParts[1], nil
}

func GetAPIKey(headers http.Header) (string, error) {
	auth := strings.TrimSpace(headers.Get("Authorization"))
	if auth == "" {
		return "", fmt.Errorf("no authorization header found")
	}

	tokenParts := strings.Split(auth, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "apikey" {
		return "", fmt.Errorf("malformed bearer token: %s", auth)
	}

	return tokenParts[1], nil
}
