package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
)

func main() {
	//http.HandleFunc("/", receiveClientRequest)
	r := mux.NewRouter()

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
		return
	}
	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}
