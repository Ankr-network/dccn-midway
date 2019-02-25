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
		return "", errors.New("DETAILNo sessionToken")
	}
	if len(sessionTokenarray) == 1 {
		return "", errors.New("DETAILInvalid sessionToken syntax")
	}

	sessionToken := sessionTokenarray[1]
	return sessionToken, nil
}


func ParseError(err error) string{
	s1 := err.Error()
	index := strings.Index(s1, "desc")
	s2 := s1[index+7:]
	//index2 := strings.Index(s2, "\"")
	//s3 := s2[:index2]
	return s2
   }
   