package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/grpc"
	"github.com/Ankr-network/dccn-midway/util"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
)
var ENDPOINT string = "client-dev.dccn.ankr.network:50051"

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

// Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Password string `json:"Password"`
	Email string `json:"Email"`
	Name     string `json:"Name"`
	Nickname string `json:"Nickname"`
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
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		return
	}

	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
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
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		log.Printf("Something went wrong! %s\n", err)
		return
	}

	log.Printf("login Successful!")
	/*JsonAuthenticationResult, err := json.Marshal(rsp.AuthenticationResult)
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		log.Printf("Something went wrong in Marshall Request! %s\n", err)
		return
	}*/
	//w.Write(JsonAuthenticationResult)
	JsonUser, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
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
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		return
	}
	log.Info(ENDPOINT)
	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		w.Write([]byte(err.Error()))
		log.Info("did not connect: ", err)
		http.Error(w, util.ParseError(err),http.StatusUnauthorized)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			http.Error(w, util.ParseError(err), http.StatusBadRequest)
			return
		}
	}(conn)
	userClient := usermgr.NewUserMgrClient(conn)
	log.Printf("Email: %s \n", Creds.Email)
	log.Printf("password: %s \n", Creds.Password)
	attribute := usermgr.UserAttributes{
		Name: Creds.Name,
	}
	_, err = userClient.Register(context.Background(), &usermgr.RegisterRequest{Password: Creds.Password, User: &usermgr.User{Email: Creds.Email, Attributes: &attribute}})
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		log.Printf("Something went wrong! \n")
		return
	}
	log.Printf("Register Success!")
	w.Write([]byte("Register Success!"))
}

/*func Welcome(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := util.SessionTokenValue(w, r)
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		return
	}
	sessionToken := c
	w.Write([]byte(fmt.Sprintf("Welcome %s!", sessionToken)))
}*/

func Refresh(w http.ResponseWriter, r *http.Request) {
	var refreshtoken RefreshToken
	err := json.NewDecoder(r.Body).Decode(&refreshtoken)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		return
	}

	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		return
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		}
	}(conn)
	userClient := usermgr.NewUserMgrClient(conn)

	rsp, err := userClient.RefreshSession(context.Background(), &usermgr.RefreshToken{RefreshToken: refreshtoken.RefreshTokenValue})
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		log.Printf("Something went wrong! \n %s", err)
		return
	}

	log.Printf("Refresh Success!")
	JsonAuthenticationResult, err := json.Marshal(rsp)
	w.Write(JsonAuthenticationResult)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	var refreshtoken RefreshToken
	err := json.NewDecoder(r.Body).Decode(&refreshtoken)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		return
	}

	c, err := util.SessionTokenValue(w, r)
	if err != nil {
		if err == errors.New("Error! No sessionToken") {
			http.Error(w, util.ParseError(err), http.StatusUnauthorized)
			return
		}
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		return
	}
	sessionToken := c

	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		return
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		}
	}(conn)

	dc := usermgr.NewUserMgrClient(conn)
	md := metadata.New(map[string]string{
		"token": sessionToken,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err = dc.Logout(ctx, &usermgr.RefreshToken{RefreshToken: refreshtoken.RefreshTokenValue})
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		log.Printf("Something went wrong! \n %s", err)
		return
	}

	log.Printf("Logout Success!")
	w.Write([]byte("Logout Success!"))
}