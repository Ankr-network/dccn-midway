package main
import (
	"net/http"
	"errors"
	"strings"
)


func sessionTokenValue(w http.ResponseWriter, r *http.Request)(string, error) {
	// We can obtain the session token from the requests cookies, which come with every request
	c := r.Header.Get("Authorization")
	/*if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return "", err
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return "", err
	}*/

	sessionTokenarray := strings.Split(c, " ")
	if sessionTokenarray == nil {
		return "", errors.New("Error! No sessionToken")
	}
	if len(sessionTokenarray) == 1 {
		return "", errors.New("Error! Invalid sessionToken syntax!")
	}
	sessionToken := sessionTokenarray[1]
	// We then get the name of the user from our cache, where we set the session token
	/*response, err := cache.Do("GET", sessionToken)
	if err != nil {
		// If there is an error fetching from cache, return an internal server error status
		w.WriteHeader(http.StatusInternalServerError)
		return "", err
	}
	if response == nil {
		// If the session token is not present in cache, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return "", err
	}*/
	// Finally, return the welcome message to the user
	return sessionToken, nil
}