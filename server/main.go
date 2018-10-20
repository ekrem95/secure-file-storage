package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ekrem95/secure-file-storage/database"
)

// Start the http server
func Start() {
	http.Handle("/", middleware(http.HandlerFunc(index)))
	http.HandleFunc("/user/login", login)
	http.HandleFunc("/user/register", register)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, string("data"))
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.URL.Query()["access_token"]
		if len(bearer) != 1 {
			errorHandler(w, r, http.StatusUnauthorized, "Missing access token")
			return
		}

		err := database.CheckAuth(bearer[0])
		if err != nil {
			errorHandler(w, r, http.StatusUnauthorized, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}
