package main

import (
	//"errors"
	//"fmt"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	//common_proto "github.com/Ankr-network/dccn-common/protos/common"
	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/grpc"
	"google.golang.org/grpc"
	"context"
)

type Confirmation struct {
	Email string `json:"Email"`
	VerificationCode string `json:"VerificationCode"`
}
type ConfirmationPassword struct {
	Email string `json:"Email"`
	VerificationCode string `json:"ConfirmationCode"`
	NewPassword string `json:"NewPassword"`
}

type ForgetUser struct {
	Email string `json:"Email"`
}

type ChangePasswordRequest struct {
     OldPassword string `json:"OldPassword"`
}

type ChangeEmailRequest struct {
	NewEmail string `json:"NewEmail"`
}

type UpdateAttributesRequest struct {
    user_id string `json:"UserId"`
    update_attributes_code string `json:"UpdateAttributeCode"`;
	name string `json:"Name"`;
    // password
    hash_password string `json:"HashPassword"`;
    tokens string `json:"Tokends"`;
    // public key of tendermint wallet
    pub_key string `json:"PubKey"`;
    creation_date uint64 `json:"CreationDate"`; //task creation date
    last_modified_date uint64 `json:"LastModifiedDate"`; //task creation date
    status string `json:"Status"`; // user's status in db
}

func confirmRegistration(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("confirm Registration")
	var Creds Confirmation
	err := json.NewDecoder(r.Body).Decode(&Creds)
	if err != nil {
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
			w.Write([]byte(err.Error()))
			log.Println(err.Error())
		}
	}(conn)

	userClient := usermgr.NewUserMgrClient(conn)
	log.Printf("Email: %s \n", Creds.Email)
	log.Printf("VerificationCode: %s \n", Creds.VerificationCode)

	_, err = userClient.ConfirmRegistration(context.TODO(), &usermgr.ConfirmRegistrationRequest{Email: Creds.Email, VerificationCode: Creds.VerificationCode})
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Something went wrong! %s\n", err)
		return
	}

	log.Printf("Request Success: %s\n", Creds.Email)
}
func forgotPassword(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Forget Password")
	log.Printf("confirm Registration")
	var Creds ForgetUser 
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&Creds)
	if err != nil {
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
			w.Write([]byte(err.Error()))
			log.Println(err.Error())
		}
	}(conn)

	userClient := usermgr.NewUserMgrClient(conn)
	log.Printf("Email: %s \n", Creds.Email)

	_, err = userClient.ForgetPassword(context.TODO(), &usermgr.ForgetPasswordRequest{Email: Creds.Email})
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Something went wrong! %s\n", err)
		return
	}
	log.Printf("Request Success: %s\n", Creds.Email)
}

func confirmPassword(w http.ResponseWriter, r *http.Request) {
	var Creds ConfirmationPassword
	err := json.NewDecoder(r.Body).Decode(&Creds)
	if err != nil {
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
			w.Write([]byte(err.Error()))
			log.Println(err.Error())
		}
	}(conn)

	userClient := usermgr.NewUserMgrClient(conn)
	log.Printf("Email: %s \n", Creds.Email)
	log.Printf("VerificationCode: %s \n", Creds.VerificationCode)
	log.Printf("NewPassword: %s \n", Creds.NewPassword)

	_, err = userClient.ConfirmPassword(context.TODO(), &usermgr.ConfirmPasswordRequest{Email: Creds.Email, VerificationCode: Creds.VerificationCode, NewPassword: Creds.NewPassword})
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Confirm Reqistration Request went wrong! %s\n", err)
		return
	}

	log.Printf("Request Success: %s\n", Creds.Email)
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Change Password")
	var Creds ChangePasswordRequest
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&Creds)
	if err != nil {
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
	log.Printf("OldPassword: %s \n", Creds.OldPassword)

	_, err = userClient.ChangePassword(context.TODO(), &usermgr.ChangePasswordRequest{OldPassword: Creds.OldPassword})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Something went wrong! %s\n", err)
		return
	}
	log.Printf("Change Password Request Success: %s\n", Creds.Email)
}


func ChangeEmail(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Change Email")
	var Creds ChangeEmailRequest
	err := json.NewDecoder(r.Body).Decode(&Creds)
	if err != nil {
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
	log.Printf("OldEmail: %s \n", Creds.OldEmail)

	_, err = userClient.ChangeEmail(context.TODO(), &usermgr.ChangeEmailRequest{NewEmail: Creds.NewEmail})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Something went wrong! %s\n", err)
		return
	}
	log.Printf("Change Email Request Success: %s\n", Creds.Email)
}


func UpdateAttribute(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Update Attribute")
	sessionToken, err:= sessionTokenValue(w, r)
	if err != nil {
		log.Info("Cannot Access to sessionToken!")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Info(sessionToken)

	var Heretask UpdateAttributesRequest
	// Get the JSON body and decode into credentials
	err1 := json.NewDecoder(r.Body).Decode(&Heretask)
	if err1 != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Something went wrong! %s\n", err1)
		return
	}

	log.Info(Heretask)
	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		log.Info("did not connect: ", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	defer conn.Close()
	dc := taskmgr.NewUserMgrClient(conn)
	md := metadata.New(map[string]string{
		"token": sessionToken,
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	task := usermgr.UserAttribute{
		Name:       	  Heretask.Name,
		HashPassword:     Heretask.HashPassword,
		Tokens:        Heretask.Tokens,
		PubKey:   Heretask.PubKey,
		CreationDate: Heretask.CreationDate,
		LastModifiedDate:	  Heretask.LastModifiedDate,
	}
	switch Heretask.Status {
	case "0":
		task.Type = usermgr.UserStatus_WAIT_ACTIVATED
	case "1":
		task.Type = usermgr.UserStatus_ACTIVATED
	case "2":
		task.Type = usermgr.UserStatus_DELETED
	default:
		task.Type = usermgr.UserStatus_WAIT_ACTIVATED
	}
	log.Info(task)
	tcrq := taskmgr.UpdateAttributesRequest{
		UserId: Heretask.UserId,
		UpdateAttributesCode: Heretask.UpdateAttributesCode,
		UserAttribute:   &task,
	}
	_, err = dc.UpdateAttributes(ctx, &tcrq)
	if err != nil {
		log.Info("Error! \n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Info("User updated successfully. \n")
}