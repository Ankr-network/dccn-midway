package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	//"time"
	"errors"
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

	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Info("did not connect: ", err)
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
	sessionToken := rsp.Token
	w.Write([]byte(fmt.Sprintf("%s", sessionToken)))
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
		log.Printf("Something went wrong! \n")
	} else {
		log.Printf("Register Success!")
	}
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := sessionTokenValue(w, r)
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
	sessionToken := c
	w.Write([]byte(fmt.Sprintf("Welcome %s!", sessionToken)))
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := sessionTokenValue(w, r)
	if err != nil {
		if err == errors.New("Error! No sessionToken") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c

	response, err := cache.Do("GET", sessionToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if response == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
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
	w.Write([]byte(fmt.Sprintf("%s", sessionToken)))
}
 