package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/grpc"
	"google.golang.org/grpc"

	//"github.com/satori/go.uuid"
	"context"

	log "github.com/sirupsen/logrus"
)

const ENDPOINT = "client-dev.dccn.ankr.network:50051"

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

// Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
}

func Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the expected password from our in memory map
	//expectedPassword, ok := users[creds.Username]

	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			log.Println(err.Error())
		}
	}(conn)
	userClient := usermgr.NewUserMgrClient(conn)
	fmt.Printf("username: %s \n", creds.Username)
	fmt.Printf("password: %s \n", creds.Password)

	rsp, err := userClient.Login(context.TODO(), &usermgr.LoginRequest{Email: creds.Username, Password: creds.Password})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Printf("Something went wrong! %s\n", err)
		return
	}

	log.Printf("login Success: %s\n", rsp.Token)
	if rsp.Error != nil {
		fmt.Printf("Something went wrong! In hub %s\n", rsp.Error)
		return
	}
	//u, err := uuid.NewV4()
	sessionToken := rsp.Token
	sessionUserid := rsp.UserId
	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of 120 seconds
	_, err = cache.Do("SETEX", sessionToken, "120", creds.Username)
	if err != nil {
		// If there is an error in setting the cache, return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds, the same as the cache
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(120 * time.Second),
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "user_id",
		Value:   sessionUserid,
		Expires: time.Now().Add(120 * time.Second),
	})
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the expected password from our in memory map
	//expectedPassword, ok := users[creds.Username]

	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			log.Println(err.Error())
		}
	}(conn)
	userClient := usermgr.NewUserMgrClient(conn)
	fmt.Printf("username: %s \n", creds.Username)
	fmt.Printf("password: %s \n", creds.Password)

	user := &usermgr.User{
		Name:     creds.Name,
		Nickname: creds.Nickname,
		Email:    creds.Username,
		Password: creds.Password,
		Balance:  0,
	}

	_, err = userClient.Register(context.Background(), user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Printf("Something went wrong! \n")
	} else {
		log.Printf("Register Success!")
	}
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	// We then get the name of the user from our cache, where we set the session token
	response, err := cache.Do("GET", sessionToken)
	if err != nil {
		// If there is an error fetching from cache, return an internal server error status
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if response == nil {
		// If the session token is not present in cache, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Finally, return the welcome message to the user
	w.Write([]byte(fmt.Sprintf("Welcome %s!", response)))
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	response, err := cache.Do("GET", sessionToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if response == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// The code uptil this point is the same as the first part of the `Welcome` route

	// Now, create a new session token for the current user

	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			log.Println(err.Error())
		}
	}(conn)
	userClient := usermgr.NewUserMgrClient(conn)

	_, err = userClient.VerifyAndRefreshToken(context.Background(), &usermgr.Token{Token: sessionToken})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Printf("Something went wrong! \n %s", err)
		return
	}

	log.Printf("Refresh Success!")
	newSessionToken := sessionToken

	_, err = cache.Do("SETEX", newSessionToken, "120", fmt.Sprintf("%s", response))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Delete the older session token
	_, err = cache.Do("DEL", sessionToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the new token as the users `session_token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   newSessionToken,
		Expires: time.Now().Add(120 * time.Second),
	})

}
