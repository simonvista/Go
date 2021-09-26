package main

import (
	"log"
	"net/http"
)


func main() {
	http.HandleFunc("/login",Login)
	http.HandleFunc("/home",Home)
	http.HandleFunc("/refresh",Refresh)
	// start server
	log.Fatal(http.ListenAndServe(":8080",nil))
}