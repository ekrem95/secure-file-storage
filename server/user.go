package server

import (
	"database/sql"
	"net/http"

	"github.com/ekrem95/secure-file-storage/database"
	"github.com/thedevsaddam/govalidator"
)

var (
	loginRules = govalidator.MapData{
		"email":    []string{"required", "min:4", "max:50", "email"},
		"password": []string{"required", "min:8"},
	}
	registerRules = govalidator.MapData{
		"name":     []string{"required", "between:3,50"},
		"email":    []string{"required", "min:4", "max:50", "email"},
		"password": []string{"required", "min:8"},
	}
)

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errorHandler(w, r, http.StatusNotFound, "")
		return
	}

	var user database.User

	e := validate(r, loginRules, &user)
	// if there is an error
	if len(e) != 0 {
		message := validationMessage(e)

		errorHandler(w, r, http.StatusBadRequest, message)
		return
	}

	res, err := database.QueryRow("SELECT id, password FROM users WHERE email = $1", user.Email)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, "")
		return
	}

	var (
		userID int
		hash   string
	)

	if err = res.Scan(&userID, &hash); err != nil {
		// If an entry with the email does not exist
		if err == sql.ErrNoRows {
			errorHandler(w, r, http.StatusBadRequest, "User does not exist")
			return
		}
		// If the error is of any other type, send a 500 status
		errorHandler(w, r, http.StatusInternalServerError, "")
		return
	}

	if match := user.CheckPasswordHash(user.Password, hash); !match {
		errorHandler(w, r, http.StatusBadRequest, "Password does not match the confirm password")
		return
	}

	var t database.Token
	database.GiveAccess(&t, userID)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errorHandler(w, r, http.StatusNotFound, "")
		return
	}

	var user database.User

	e := validate(r, registerRules, &user)
	// if there is an error
	if len(e) != 0 {
		message := validationMessage(e)

		errorHandler(w, r, http.StatusBadRequest, message)
		return
	}

	errorHandler(w, r, http.StatusBadRequest, "message")
}
