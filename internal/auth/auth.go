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
    if len(vals) != 1 {
        return "", errors.New("malformed authentication header")
    }

    return vals[0], nil
}
