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

type RefreshToken struct{
	RefreshTokenValue string `json:"RefreshToken"`
}


func Signin(w http.ResponseWriter, r *http.Request) {
	var NewUser Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&NewUser)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusUnauthorized)
		log.Info("did not connect: ", err)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			log.Println(err.Error())
		}
	}(conn)
	userClient := usermgr.NewUserMgrClient(conn)
	log.Printf("Email: %s \n", NewUser.Email)
	log.Printf("password: %s \n", NewUser.Password)

	rsp, err := userClient.Login(context.TODO(), &usermgr.LoginRequest{Email: NewUser.Email, Password: NewUser.Password})
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("Something went wrong! %s\n", err)
		return
	}

	log.Printf("login Successful!")
	w.Write([]byte("login Successful!"))
	JsonAuthenticationResult, err := json.Marshal(rsp.AuthenticationResult)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Something went wrong in Marshall Request! %s\n", err)
		return
	}
	w.Write(JsonAuthenticationResult)
	JsonUser, _err := json.Marshal(rsp.User)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Something went wrong in Marshall Request! %s\n", err)
		return
	}
	w.Write(JsonUser)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var Creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&Creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		w.Write([]byte(err.Error()))
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
	log.Printf("Email: %s \n", Creds.Email)
	log.Printf("password: %s \n", Creds.Password)
	attribute := usermgr.UserAttribute{
		Name: Creds.Name,
	}
	_, err = userClient.Register(context.Background(), &usermgr.RegisterRequest{Password: Creds.Password, User: &usermgr.User{Email: Creds.Email, Attribute: &attribute}})
	if err != nil {
		w.Write([]byte(err.Error()))
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
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c
	w.Write([]byte(fmt.Sprintf("Welcome %s!", sessionToken)))
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	var refreshtoken RefreshToken
	err := json.NewDecoder(r.Body).Decode(&refreshtoken)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.Write([]byte(err.Error()))
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

	_, err = userClient.RefreshSession(context.Background(), &usermgr.RefreshToken{RefreshToken: refreshtoken.RefreshTokenValue})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("Something went wrong! \n %s", err)
		return
	}

	log.Printf("Refresh Success!")
}