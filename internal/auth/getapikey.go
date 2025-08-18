package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	value, ok := headers["Authorization"]
	if !ok {
		log.Println("Error 1")
		return "", errors.New("no authorization field in header")
	} else if len(value) == 0 {
		log.Println("Error 2")
		return "", errors.New("authorization field is empty")
	}

	auth := value[0]
	if len(auth) == 0 {
		log.Println("Error 3")
		return "", errors.New("no token provided")
	}

	if !strings.HasPrefix(auth, "ApiKey ") {
		log.Println("Error 4")
		return "", errors.New("no token bearer")
	}

	trimed := strings.TrimPrefix(auth, "ApiKey ")
	if len(trimed) < 1 {
		log.Println("Error 5")
		return "", errors.New("no token provided")
	}

	fullTrimed := strings.TrimSpace(trimed)
	return fullTrimed, nil
}
