package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

const (
	ADDRESS = "http://localhost:8000"
)

var (
	urlLogin  = fmt.Sprintf("%s/login", ADDRESS)
	urlCreate = fmt.Sprintf("%s/create", ADDRESS)
	urlUpdate = fmt.Sprintf("%s/update", ADDRESS)
	urlPurge  = fmt.Sprintf("%s/purge", ADDRESS)
	urlCancel = fmt.Sprintf("%s/delete", ADDRESS)
	urlList   = fmt.Sprintf("%s/list", ADDRESS)
	urlSignup = fmt.Sprintf("%s/signup", ADDRESS)
)

type Person struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Age     int    `json:"age"`
}

func TestSignin(t *testing.T) {
	t.Log("URL for login:>", urlLogin)
	client := &http.Client{}

	var jsonStrlogin = []byte(`{"username":"xiaowu", "password":"1111"}`)
	reqlogin, err := http.NewRequest("POST", urlLogin, bytes.NewBuffer(jsonStrlogin))
	reqlogin.Header.Set("X-Custom-Header", "myvalue")
	reqlogin.Header.Set("Content-Type", "application/json")

	resplogin, err := client.Do(reqlogin)
	if err != nil {
		//panic(err)
		t.Error(err)
	}
	defer resplogin.Body.Close()

	t.Log("Signin response Status:", resplogin.Status)
	if resplogin.Status != "200 OK" {
		t.Error("login Status Error!")
	}

	//fmt.Println("Signin response Headers:", resplogin.Header)
	body, _ := ioutil.ReadAll(resplogin.Body)
	t.Log("Signin response Body:", string(body))
}

func TestSigninBAD(t *testing.T) {
	t.Log("URL for login:>", urlLogin)
	client := &http.Client{}

	var jsonStrlogin = []byte(`{"username":"xiaowu", "password":"11111"}`)
	reqlogin, err := http.NewRequest("POST", urlLogin, bytes.NewBuffer(jsonStrlogin))
	reqlogin.Header.Set("X-Custom-Header", "myvalue")
	reqlogin.Header.Set("Content-Type", "application/json")

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

func TestSignup(t *testing.T) {
	urlsignup := "http://localhost:8000/signup"
	t.Log("URL for signup:>", urlsignup)
	var jsonStrSignup = []byte(`{"username":"xiaowu", "password":"1111"}`)
	reqsignup, err := http.NewRequest("POST", urlsignup, bytes.NewBuffer(jsonStrSignup))
	reqsignup.Header.Set("X-Custom-Header", "myvalue")
	reqsignup.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	respsignup, err := client.Do(reqsignup)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respsignup.Body.Close()

	t.Log("Signup response Status:", respsignup.Status)
	if respsignup.Status != "200 OK" {
		t.Error("Signup Status Error!")
	}
	t.Log("Signup response Headers:", respsignup.Header)
	body, _ := ioutil.ReadAll(respsignup.Body)
	t.Log("Signup response Body:", string(body))
}

func TestSignupBad(t *testing.T) {
	t.Log("URL for signup:>", urlSignup)
	var jsonStrSignup = []byte(`{"username":"xiaowu", "password":"11111"}`)
	reqsignup, err := http.NewRequest("POST", urlSignup, bytes.NewBuffer(jsonStrSignup))
	reqsignup.Header.Set("X-Custom-Header", "myvalue")
	reqsignup.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	respsignup, err := client.Do(reqsignup)
	if err != nil {
		//panic(err)
		t.Error("err")
	}
	defer respsignup.Body.Close()

	t.Log("Signup response Status:", respsignup.Status)
	if respsignup.Status == "200 OK" {
		t.Error("Signup Status Repeat Error!")
	}
	t.Log("Signup response Headers:", respsignup.Header)
	body, _ := ioutil.ReadAll(respsignup.Body)
	t.Log("Signup response Body:", string(body))
}
