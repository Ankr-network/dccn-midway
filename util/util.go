package util

import (
	"errors"
	"net/http"
	"strings"
)

func SessionTokenValue(w http.ResponseWriter, r *http.Request) (string, error) {
	// We can obtain the session token from the requests cookies, which come with every request
	c := r.Header.Get("Authorization")

	sessionTokenarray := strings.Split(c, " ")
	if sessionTokenarray == nil {
		return "", errors.New("No sessionToken")
	}
	if len(sessionTokenarray) == 1 {
		return "", errors.New("Invalid sessionToken syntax")
	}

	sessionToken := sessionTokenarray[1]
	return sessionToken, nil
}
