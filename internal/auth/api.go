package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	tokenSlice := strings.Fields(authHeader)
	if len(tokenSlice) != 2 {
		return "", errors.New("wrong token format")
	}
	if strings.ToLower(tokenSlice[0]) != "apikey" {
		return "", errors.New("wrong header")
	}
	tokenStr := tokenSlice[1]

	return tokenStr, nil
}
