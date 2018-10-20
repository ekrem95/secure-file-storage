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

type successResponse struct {
	Atoken string `json:"access_token"`
}

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

	result, err := database.QueryRow("SELECT id, password FROM users WHERE email = $1", user.Email)
	if err != nil {
		internalError(w, r)
		return
	}

	var (
		userID int
		hash   string
	)

	if err = result.Scan(&userID, &hash); err != nil {
		if err == sql.ErrNoRows {
			errorHandler(w, r, http.StatusBadRequest, "User does not exist")
			return
		}
		// If the error is of any other type, send internalError
		internalError(w, r)
		return
	}

	if match := user.CheckPasswordHash(user.Password, hash); !match {
		errorHandler(w, r, http.StatusBadRequest, "Password does not match the confirm password")
		return
	}

	var t database.Token
	ss, err := database.NewJWTWithClaims(&t, userID)
	if err != nil {
		internalError(w, r)
		return
	}

	res := successResponse{Atoken: ss}

	response(w, r, res)
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

	result, err := database.QueryRow("SELECT id FROM users WHERE email = $1", user.Email)
	if err != nil {
		internalError(w, r)
		return
	}

	err = result.Scan(&user.ID)
	switch err {
	case sql.ErrNoRows:
		err = user.Save()
		if err != nil {
			internalError(w, r)
			return
		}

		var t database.Token
		ss, err := database.NewJWTWithClaims(&t, user.ID)
		if err != nil {
			internalError(w, r)
			return
		}

		res := successResponse{Atoken: ss}

		response(w, r, res)
	case nil:
		errorHandler(w, r, http.StatusBadRequest, "An account for the specified email address already exists")
	default:
		// If the error is of any other type, send internalError
		internalError(w, r)
	}
}
