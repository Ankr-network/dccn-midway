package main

import (
	"log"
	"net/http"

	"github.com/Ankr-network/dccn-midway/handlers"
	"github.com/gorilla/mux"
)


func main() {
	r := mux.NewRouter()

	r.HandleFunc("/login", handlers.Signin)
	r.HandleFunc("/signup", handlers.Signup)
//	r.HandleFunc("/welcome", handlers.Welcome)
	r.HandleFunc("/refresh", handlers.Refresh)
	r.HandleFunc("/logout", handlers.Logout)
	r.HandleFunc("/create", handlers.CreateTask)
	r.HandleFunc("/update", handlers.UpdateTask)
	r.HandleFunc("/list", handlers.ListTask)
	r.HandleFunc("/delete", handlers.CancelTask)
	r.HandleFunc("/purge", handlers.PurgeTask)
	r.HandleFunc("/dclist", handlers.DataCenterList)
	r.HandleFunc("/confirmregistration", handlers.ConfirmRegistration)
	r.HandleFunc("/forgotpassword", handlers.ForgotPassword)
	r.HandleFunc("/confirmpassword", handlers.ConfirmPassword)
	r.HandleFunc("/changepassword", handlers.ChangePassword)
	r.HandleFunc("/changeemail", handlers.ChangeEmail)
	r.HandleFunc("/refresh", handlers.Refresh)
	r.HandleFunc("/updateattribute", handlers.UpdateAttribute)
	r.HandleFunc("/taskoverview", handlers.TaskOverview)
	r.HandleFunc("/taskleaderboard", handlers.TaskLeaderBoard)
	r.HandleFunc("/networkinfo", handlers.NetworkInfo)
	r.HandleFunc("/dcleaderboard", handlers.DCLeaderBoard)
	//http.HandleFunc("/confirmregistration", handlers.confirmRegistration)
	//http.HandleFunc("/forgotpassword", handlers.forgotPassword)
	//http.HandleFunc("/confirmpassword", handlers.confirmPassword)
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
		return
	}
	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}
