package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	value, ok := headers["Authorization"]
	if !ok {
		return "", errors.New("no authorization field in header")
	} else if len(value) == 0 {
		return "", errors.New("authorization field is empty")
	}

	auth := value[0]
	if len(auth) == 0 {
		return "", errors.New("no token provided")
	}

	if !strings.HasPrefix(auth, "Bearer ") {
		return "", errors.New("no token bearer")
	}

	trimed := strings.TrimPrefix(auth, "Bearer ")
	if len(trimed) < 1 {
		return "", errors.New("no token provided")
	}

	fullTrimed := strings.TrimSpace(trimed)
	return fullTrimed, nil
}
