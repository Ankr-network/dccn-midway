package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/grpc"
	"github.com/Ankr-network/dccn-midway/util"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var ENDPOINT string

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
	Creds := usermgr.User{}
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&Creds)
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
	fmt.Printf("Email: %s \n", Creds.Email)
	fmt.Printf("password: %s \n", Creds.Password)

	rsp, err := userClient.Login(context.TODO(), &usermgr.LoginRequest{Email: Creds.Email, Password: Creds.Password})
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
	Creds := usermgr.User{}
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&Creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		log.Info("did not connect: ", err)
		w.WriteHeader(http.StatusUnauthorized)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}(conn)
	userClient := usermgr.NewUserMgrClient(conn)
	fmt.Printf("Email: %s \n", Creds.Email)
	fmt.Printf("password: %s \n", Creds.Password)

	_, err = userClient.Register(context.Background(), &Creds)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Info(err)
		log.Printf("Something went wrong! \n")
	} else {
		log.Printf("Register Success!")
	}
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := util.SessionTokenValue(w, r)
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
	c, err := util.SessionTokenValue(w, r)
	if err != nil {
		if err == errors.New("Error! No sessionToken") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c

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
