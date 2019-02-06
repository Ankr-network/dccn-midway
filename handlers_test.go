
package main
 
import (
	"testing"
	"net/http"
	"fmt"
	"io/ioutil"
	"bytes"
)
 
const (
	ADDRESS = "localhost:8080"
)
 
type Person struct {
	Name string `json:"name"`
	Address string `json:"address"`
	Age int `json:"age"`
}
 
func TestSignin(t *testing.T) {
	urllogin := "http://localhost:8000/login"
    fmt.Println("URL for login:>", urllogin)
	client := &http.Client{}

	var jsonStrlogin = []byte(`{"username":"xiaowu", "password":"1111"}`)
    reqlogin, err := http.NewRequest("POST", urllogin, bytes.NewBuffer(jsonStrlogin))
    reqlogin.Header.Set("X-Custom-Header", "myvalue")
    reqlogin.Header.Set("Content-Type", "application/json")

    resplogin, err := client.Do(reqlogin)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
    defer resplogin.Body.Close()

	fmt.Println("Signin response Status:", resplogin.Status)
	if resplogin.Status != "200 OK"{
		t.Error("login Status Error!")
	}

    //fmt.Println("Signin response Headers:", resplogin.Header)
    body, _ := ioutil.ReadAll(resplogin.Body)
    fmt.Println("Signin response Body:", string(body))
}

func TestSigninBAD(t *testing.T) {
	urllogin := "http://localhost:8000/login"
    fmt.Println("URL for login:>", urllogin)
	client := &http.Client{}

	var jsonStrlogin = []byte(`{"username":"xiaowu", "password":"11111"}`)
    reqlogin, err := http.NewRequest("POST", urllogin, bytes.NewBuffer(jsonStrlogin))
    reqlogin.Header.Set("X-Custom-Header", "myvalue")
    reqlogin.Header.Set("Content-Type", "application/json")

    resplogin, err := client.Do(reqlogin)
    if err != nil {
		//panic(err)
		t.Error("err")
    }
    defer resplogin.Body.Close()

	fmt.Println("Signin response Status:", resplogin.Status)
	if resplogin.Status == "200 OK"{
		t.Error("login Status Error!")
	}

    //fmt.Println("Signin response Headers:", resplogin.Header)
    body, _ := ioutil.ReadAll(resplogin.Body)
    fmt.Println("Signin response Body:", string(body))
}


func TestSignup(t *testing.T) {
	urlsignup := "http://localhost:8000/signup"
	fmt.Println("URL for signup:>", urlsignup)
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

	fmt.Println("Signup response Status:", respsignup.Status)
	if respsignup.Status != "200 OK"{
		t.Error("Signup Status Error!")
	}
    fmt.Println("Signup response Headers:", respsignup.Header)
    body, _ := ioutil.ReadAll(respsignup.Body)
    fmt.Println("Signup response Body:", string(body))
}

func TestSignupBad(t *testing.T) {
	urlsignup := "http://localhost:8000/signup"
	fmt.Println("URL for signup:>", urlsignup)
    var jsonStrSignup = []byte(`{"username":"xiaowu", "password":"11111"}`)
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

	fmt.Println("Signup response Status:", respsignup.Status)
	if respsignup.Status == "200 OK"{
		t.Error("Signup Status Repeat Error!")
	}
    fmt.Println("Signup response Headers:", respsignup.Header)
    body, _ := ioutil.ReadAll(respsignup.Body)
    fmt.Println("Signup response Body:", string(body))
}