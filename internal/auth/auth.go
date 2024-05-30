package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
    authHeader := headers.Get("Authorization")
    if authHeader == "" {
        return "", errors.New("no authentication info found")
    }
    
    vals := strings.Split(authHeader, " ")
    if len(vals) != 2 {
        return "", errors.New("malformed authentication header")
    }

    if vals[0] != "ApiKey" {
        return "", errors.New("malformed first part of authentication header")
    }

    return vals[1], nil
}
