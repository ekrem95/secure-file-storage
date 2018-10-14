package server

import (
	"log"
	"net/http"
)

// Start the http server
func Start() {
	http.HandleFunc("/user/login", login)
	http.HandleFunc("/user/register", login)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
