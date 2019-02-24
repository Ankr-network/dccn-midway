package main

import (
	"log"
	"net/http"

	"github.com/Ankr-network/dccn-midway/handlers"
	"github.com/gorilla/mux"
)


func main() {
	r := mux.NewRouter()

	// user management
	r.HandleFunc("/signup", handlers.Signup) // POST
	r.HandleFunc("/confirm_registration", handlers.ConfirmRegistration) // POST
	r.HandleFunc("/login", handlers.Signin) // POST
	r.HandleFunc("/logout", handlers.Logout) // POST
	r.HandleFunc("/refresh", handlers.Refresh) // POST
	r.HandleFunc("/forgot_password", handlers.ForgotPassword) // POST
	r.HandleFunc("/confirm_password", handlers.ConfirmPassword) // POST
	r.HandleFunc("/change_password", handlers.ChangePassword) // POST
	r.HandleFunc("/change_email", handlers.ChangeEmail) // POST
	r.HandleFunc("/update_attribute", handlers.UpdateAttribute) // POST
	// r.HandleFunc("/welcome", handlers.Welcome)
	
	// task management
	r.HandleFunc("/task/create", handlers.CreateTask) // POST
	r.HandleFunc("/task/update", handlers.UpdateTask) // POST
	r.HandleFunc("/task/list", handlers.ListTask) // GET
	r.HandleFunc("/task/delete", handlers.CancelTask) // POST
	r.HandleFunc("/task/purge", handlers.PurgeTask) // POST

	// data center management
	r.HandleFunc("/dc/list", handlers.DataCenterList) // GET
	r.HandleFunc("/dc/task_overview", handlers.TaskOverview) // GET
	r.HandleFunc("/dc/task_leaderboard", handlers.TaskLeaderBoard) // GET
	r.HandleFunc("/dc/network_info", handlers.NetworkInfo) // GET
	r.HandleFunc("/dc/leaderboard", handlers.DCLeaderBoard) // GET

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
