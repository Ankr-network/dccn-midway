package handlers

import (
	"io/ioutil"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	//common_proto "github.com/Ankr-network/dccn-common/protos/common"
	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/grpc"
	"google.golang.org/grpc"
	"context"
	"github.com/Ankr-network/dccn-midway/util"
	metadata "google.golang.org/grpc/metadata"
)

type Confirmation struct {
	Email string `json:"Email"`
	ConfirmationCode string `json:"ConfirmationCode"`
}
type ConfirmationPassword struct {
	Email string `json:"Email"`
	ConfirmationCode string `json:"ConfirmationCode"`
	NewPassword string `json:"NewPassword"`
}

type ForgetUser struct {
	Email string `json:"Email"`
}

type ChangePasswordRequest struct {
	 OldPassword string `json:"OldPassword"`
	 NewPassword string `json:"NewPassword"`
}

type ChangeEmailRequest struct {
	NewEmail string `json:"NewEmail"`
}

type UpdateAttributesRequest struct {
    user_id string `json:"UserId"`
    update_attributes_code string `json:"UpdateAttributeCode"`;
	name string `json:"Name"`;
    hash_password string `json:"HashPassword"`;
    tokens string `json:"Tokends"`;
    pub_key string `json:"PubKey"`;
    creation_date uint64 `json:"CreationDate"`; //task creation date
    last_modified_date uint64 `json:"LastModifiedDate"`; //task creation date
    status string `json:"Status"`; // user's status in db
}

func ConfirmRegistration(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("confirm Registration")
	var Creds Confirmation
	err := json.NewDecoder(r.Body).Decode(&Creds)
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		return
	}

	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		log.Info("did not connect: ", err)
		return 
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			http.Error(w, util.ParseError(err), http.StatusBadRequest)
		}
	}(conn)

	userClient := usermgr.NewUserMgrClient(conn)
	log.Printf("Email: %s \n", Creds.Email)
	log.Printf("ConfirmationCode: %s \n", Creds.ConfirmationCode)

	_, err = userClient.ConfirmRegistration(context.TODO(), &usermgr.ConfirmRegistrationRequest{Email: Creds.Email, ConfirmationCode: Creds.ConfirmationCode})
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		log.Printf("Something went wrong! %s\n", err)
		return
	}
	w.Write([]byte("Request Success!"))
	log.Printf("Request Success: %s\n", Creds.Email)
}
func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Forget Password")
	var Creds ForgetUser 
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&Creds)
	if err != nil {
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
			http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		}
	}(conn)

	userClient := usermgr.NewUserMgrClient(conn)
	log.Printf("Email: %s \n", Creds.Email)

	_, err = userClient.ForgotPassword(context.TODO(), &usermgr.ForgotPasswordRequest{Email: Creds.Email})
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		log.Printf("Something went wrong! %s\n", err)
		return
	}
	w.Write([]byte("Request Success!"))
	log.Printf("Request Success: %s\n", Creds.Email)
}

func ConfirmPassword(w http.ResponseWriter, r *http.Request) {
	var Creds ConfirmationPassword
	err := json.NewDecoder(r.Body).Decode(&Creds)
	if err != nil {
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
			http.Error(w, util.ParseError(err), http.StatusBadRequest)
		}
	}(conn)

	userClient := usermgr.NewUserMgrClient(conn)
	log.Printf("Email: %s \n", Creds.Email)
	log.Printf("ConfirmationCode: %s \n", Creds.ConfirmationCode)
	log.Printf("NewPassword: %s \n", Creds.NewPassword)

	_, err = userClient.ConfirmPassword(context.TODO(), &usermgr.ConfirmPasswordRequest{Email: Creds.Email, ConfirmationCode: Creds.ConfirmationCode, NewPassword: Creds.NewPassword})
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		log.Printf("Confirm Reqistration Request went wrong! %s\n", err)
		return
	}

	log.Printf("Request Success: %s\n", Creds.Email)
	w.Write([]byte("Confirm Success!"))
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Change Password")
	var Creds ChangePasswordRequest
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&Creds)
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		return
	}
	sessionToken, err:= util.SessionTokenValue(w, r)
	if err != nil {
		log.Info("Cannot Access to sessionToken!")
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		return
	}
	log.Info(sessionToken)

	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		log.Info("did not connect: ", err)
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

	log.Printf("OldPassword: %s \n", Creds.OldPassword)

	_, err = dc.ChangePassword(ctx, &usermgr.ChangePasswordRequest{OldPassword: Creds.OldPassword, NewPassword: Creds.NewPassword})
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		log.Printf("Something went wrong! %s\n", err)
		return
	}
	log.Printf("Change Password Request Success!")
	w.Write([]byte("Change Password Request Success!"))
}


func ChangeEmail(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Change Email")
	var Creds ChangeEmailRequest
	err := json.NewDecoder(r.Body).Decode(&Creds)
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		return
	}
	sessionToken, err:= util.SessionTokenValue(w, r)
	if err != nil {
		log.Info("Cannot Access to sessionToken!")
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		return
	}
	log.Info(sessionToken)
	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		log.Info("did not connect: ", err)
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
	log.Printf("NewEmail: %s \n", Creds.NewEmail)

	_, err = dc.ChangeEmail(ctx, &usermgr.ChangeEmailRequest{NewEmail: Creds.NewEmail})
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		log.Printf("Something went wrong! %s\n", err)
		return
	}
	log.Printf("Change Email Request Success!")
	w.Write([]byte("Change Email Request Success!"))
}

func UpdateAttribute(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	log.Printf("Update Attribute")
	sessionToken, err:= util.SessionTokenValue(w, r)
	if err != nil {
		log.Info("Cannot Access to sessionToken!")
		http.Error(w, util.ParseError(err), http.StatusUnauthorized)
		return
	}
	log.Info(sessionToken)

	conn, err := grpc.Dial(ENDPOINT, grpc.WithInsecure())
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		log.Info("did not connect: ", err)
	}

	defer conn.Close()
	dc := usermgr.NewUserMgrClient(conn)
	md := metadata.New(map[string]string{
		"token": sessionToken,
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	var Heretask []interface{}
	data, err := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(data, &Heretask); err != nil {
		log.Info(err)
        return
	}
	taskarray := []*usermgr.UserAttribute{}
	task := &usermgr.UserAttribute{}
	var temp map[string]interface{}
	for i := range Heretask {
		temp = Heretask[i].(map[string]interface{})
		switch temp["Value"].(type){
		case string:
			log.Info("string")
			data := temp["Value"]
			if str, ok := data.(string); ok {
				task = new(usermgr.UserAttribute)
				task.Key = temp["Key"].(string)
				task.Value = &usermgr.UserAttribute_StringValue{StringValue: str}
			} else {
				return
			}
		case int:
			log.Info("Int")
			data := temp["Value"]
			if str, ok := data.(int64); ok {
				task = new(usermgr.UserAttribute)
				task.Key = temp["Key"].(string)
				task.Value = &usermgr.UserAttribute_IntValue{IntValue: str}
			} else {
				return
			}
		case float64:
			log.Info("Float")
			data := temp["Value"]
			if str, ok := data.(float64); ok {
				task = new(usermgr.UserAttribute)
				task.Key = temp["Key"].(string)
				task.Value = &usermgr.UserAttribute_DoubleValue{DoubleValue: str}
			} else {
				return
			}
		case bool:
			log.Info("bool")
			data := temp["Value"]
			if str, ok := data.(bool); ok {
				task = new(usermgr.UserAttribute)
			
				task.Key = temp["Key"].(string)
				task.Value = &usermgr.UserAttribute_BoolValue{BoolValue: str}
			
			} else {
				return
			}
			
		}
		taskarray = append(taskarray, task)
	}
	tcrq := usermgr.UpdateAttributesRequest{
		UserAttributes:   taskarray,
	}
	log.Info(taskarray)
	rsp, err := dc.UpdateAttributes(ctx, &tcrq)
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		log.Info("Error! \n", err)
		return
	}
	log.Info("User updated successfully. \n")
	jsonrsp, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, util.ParseError(err), http.StatusBadRequest)
		log.Info("Marshal Error! \n", err)
		return
	}
	w.Write(jsonrsp)
}