package main

import (
	"bytes"
//	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestconfirmRegistration(t *testing.T) {
	t.Log("URL for login:>", urlLogin)
	client := &http.Client{}

	var jsonStrlogin = []byte(`{"email":"testuser@mailinator.com", "confirmation code":"confirmation code"}`)
	reqlogin, err := http.NewRequest("POST", urlLogin, bytes.NewBuffer(jsonStrlogin))

	resplogin, err := client.Do(reqlogin)
	if err != nil {
		t.Error(err)
	}
	defer resplogin.Body.Close()

	t.Log("Signin response Status:", resplogin.Status)
	if resplogin.Status != "200 OK" {
		t.Error("Registration Status Error!")
	}

	//fmt.Println("Signin response Headers:", resplogin.Header)
	body, _ := ioutil.ReadAll(resplogin.Body)
	t.Log("Signin response Body:", string(body))
}

func TestconfirmRegistrationWrongCode(t *testing.T) {
	t.Log("URL for login:>", urlLogin)
	client := &http.Client{}

	var jsonStrlogin = []byte(`{"email":"testuser@mailinator.com", "confirmation code":"wrong confirmation code"}`)
	reqlogin, err := http.NewRequest("POST", urlLogin, bytes.NewBuffer(jsonStrlogin))

	resplogin, err := client.Do(reqlogin)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer resplogin.Body.Close()

	t.Log("Signin response Status:", resplogin.Status)
	if resplogin.Status == "200 OK" {
		t.Error("login Status Error!")
	}

	//fmt.Println("Signin response Headers:", resplogin.Header)
	body, _ := ioutil.ReadAll(resplogin.Body)
	t.Log("Signin response Body:", string(body))
}

func TestconfirmRegistrationNoEmail(t *testing.T) {
	t.Log("URL for login:>", urlLogin)
	client := &http.Client{}

	var jsonStrlogin = []byte(`{"email":"wronguser@mailinator.com", "confirmation code":"confirmation code"}`)
	reqlogin, err := http.NewRequest("POST", urlLogin, bytes.NewBuffer(jsonStrlogin))

	resplogin, err := client.Do(reqlogin)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer resplogin.Body.Close()

	t.Log("Signin response Status:", resplogin.Status)
	if resplogin.Status == "200 OK" {
		t.Error("login Status Error!")
	}

	//fmt.Println("Signin response Headers:", resplogin.Header)
	body, _ := ioutil.ReadAll(resplogin.Body)
	t.Log("Signin response Body:", string(body))
}

func TestforgotPassword(t *testing.T) {
	t.Log("URL for ForgetPasswork:>", urlForgetPassword)
	var jsonStrForgetPassword = []byte(`{"email":"testuser2@mailinator.com"}`)
	reqForgetPassword, err := http.NewRequest("POST", urlForgetPassword, bytes.NewBuffer(jsonStrForgetPassword))

	client := &http.Client{}
	respForgetPassword, err := client.Do(reqForgetPassword)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respForgetPassword.Body.Close()

	t.Log("Forget response Status:", respForgetPassword.Status)
	if respForgetPassword.Status != "200 OK" {
		t.Error("Forget Status Error!")
	}
	t.Log("Forget response Headers:", respForgetPassword.Header)
	body, _ := ioutil.ReadAll(respForgetPassword.Body)
	t.Log("Forget response Body:", string(body))
}


func TestforgotPasswordNoEmail(t *testing.T) {
	t.Log("URL for Forget:>", urlForgetPassword)
	var jsonStrForgetPassword = []byte(`{"email":"wrongtestuser@mailinator.com}`)
	reqForgetPassword, err := http.NewRequest("POST", urlForgetPassword, bytes.NewBuffer(jsonStrForgetPassword))

	client := &http.Client{}
	respForgetPassword, err := client.Do(reqForgetPassword)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respForgetPassword.Body.Close()

	t.Log("Forget response Status:", respForgetPassword.Status)
	if respForgetPassword.Status == "200 OK" {
		t.Error("Forget Status Error! Cannot find an unregistered account")
	}
	t.Log("Forget response Headers:", respForgetPassword.Header)
	body, _ := ioutil.ReadAll(respForgetPassword.Body)
	t.Log("Forget response Body:", string(body))
}

func TestconfirmPassword(t *testing.T) {
	t.Log("URL for confirm Password:>", urlConfirmPassword)
	var jsonStrConfirmPassword = []byte(`{"email":"testuser2@mailinator.com", "VerificationCode": "VerificationCode", "NewPassword": "NewPassword"}`)
	reqConfirmPassword, err := http.NewRequest("POST", urlConfirmPassword, bytes.NewBuffer(jsonStrConfirmPassword))

	client := &http.Client{}
	respConfirmPassword, err := client.Do(reqConfirmPassword)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respConfirmPassword.Body.Close()

	t.Log("ConfirmPassword response Status:", respConfirmPassword.Status)
	if respConfirmPassword.Status != "200 OK" {
		t.Error("Confirm Password Status Error!")
	}
	//t.Log("Signup response Headers:", respsignup.Header)
	//body, _ := ioutil.ReadAll(respsignup.Body)
	//t.Log("Signup response Body:", string(body))
}

func TestconfirmPasswordRepeatPassword(t *testing.T) {//New Password cannot be the same with the old one
	t.Log("URL for confirm Password:>", urlConfirmPassword)
	var jsonStrConfirmPassword = []byte(`{"email":"testuser2@mailinator.com", "VerificationCode": "VerificationCode", "NewPassword": "NewPassword"}`)
	reqConfirmPassword, err := http.NewRequest("POST", urlConfirmPassword, bytes.NewBuffer(jsonStrConfirmPassword))

	client := &http.Client{}
	respConfirmPassword, err := client.Do(reqConfirmPassword)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respConfirmPassword.Body.Close()

	t.Log("ConfirmPassword response Status:", respConfirmPassword.Status)
	if respConfirmPassword.Status == "200 OK" {
		t.Error("Confirm Password Status Error! New Password should not be the same with the old one")
	}
//	t.Log("Signup response Headers:", respsignup.Header)
//	body, _ := ioutil.ReadAll(respsignup.Body)
//	t.Log("Signup response Body:", string(body))
}
func TestconfirmPasswordWrongVerificationCode(t *testing.T) {//New Password cannot be the same with the old one
	t.Log("URL for confirm Password:>", urlConfirmPassword)
	var jsonStrConfirmPassword = []byte(`{"email":"testuser2@mailinator.com", "VerificationCode": "WrongVerificationCode", "NewPassword": "NewPassword"}`)
	reqConfirmPassword, err := http.NewRequest("POST", urlConfirmPassword, bytes.NewBuffer(jsonStrConfirmPassword))

	client := &http.Client{}
	respConfirmPassword, err := client.Do(reqConfirmPassword)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respConfirmPassword.Body.Close()

	t.Log("ConfirmPassword response Status:", respConfirmPassword.Status)
	if respConfirmPassword.Status == "200 OK" {
		t.Error("Confirm Password Status Error! New Password should not be set when wrong VerificationCode is provided")
	}
//	t.Log("Confirm Password response Headers:", respsignup.Header)
//	body, _ := ioutil.ReadAll(respsignup.Body)
//	t.Log("Confirm Password response Body:", string(body))
}