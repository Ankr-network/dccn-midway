package main

import (
	"log"
	"net/http"
	"github.com/gomodule/redigo/redis"
)

var cache redis.Conn

func main() {
	initCache()
	// "Signin" and "Signup" are handler that we will implement
	http.HandleFunc("/login", Signin)
	http.HandleFunc("/signup", Signup)
	http.HandleFunc("/welcome", Welcome)
	http.HandleFunc("/refresh", Refresh)
	http.HandleFunc("/create", CreateTask)
	http.HandleFunc("/update", UpdateTask)
	http.HandleFunc("/list", ListTask)
	http.HandleFunc("/delete", CancelTask)
	http.HandleFunc("/purge", PurgeTask)
	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func initCache() {
	conn, err := redis.DialURL("redis://localhost")
	if err != nil {
		panic(err)
	}
	cache = conn
}
