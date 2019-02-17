package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gomodule/redigo/redis"
)

var cache redis.Conn

func main() {
//	initCache()
	// "Signin" and "Signup" are handler that we will implement
	//http.HandleFunc("/", receiveClientRequest)
	r := mux.NewRouter()
/*	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")             
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization") 
		w.Header().Set("content-type", "application/json")     
		log.Println("Hello")       
		return
})*///log.Println("xiaowang")
	r.HandleFunc("/login", Signin)
	r.HandleFunc("/signup", Signup)
	r.HandleFunc("/welcome", Welcome)
	r.HandleFunc("/refresh", Refresh)
	r.HandleFunc("/create", CreateTask)
	r.HandleFunc("/update", UpdateTask)
	r.HandleFunc("/list", ListTask)
	r.HandleFunc("/delete", CancelTask)
	r.HandleFunc("/purge", PurgeTask)
	r.HandleFunc("/dclist", DataCenterList)
	r.HandleFunc("/taskdetail", TaskDetail)
	//http.HandleFunc("/confirmregistration", confirmRegistration)
	//http.HandleFunc("/forgotpassword", forgotPassword)
	//http.HandleFunc("/confirmpassword", confirmPassword)
	// start the server on port 8000
	http.Handle("/", &MyServer{r})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type MyServer struct {
    r *mux.Router
}

func (s *MyServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
    if origin := req.Header.Get("Origin"); origin != "" {
	    rw.Header().Add("Access-Control-Allow-Origin", origin)
        rw.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
        rw.Header().Add("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Authorization")
    }
    // Stop here if its Preflighted OPTIONS request
    if req.Method == "OPTIONS" {
		log.Println("xiaohua")
        return
    }
    // Lets Gorilla work
    s.r.ServeHTTP(rw, req)
}