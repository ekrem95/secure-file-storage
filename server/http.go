package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ekrem95/secure-file-storage/database"
	"github.com/thedevsaddam/govalidator"
)

// ErrorResponse is a custom response type
type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
	Message    string `json:"message"`
}

const validation = "validation"

var errors = map[int]string{
	http.StatusBadRequest:          "Bad Request",
	http.StatusUnauthorized:        "Unauthorized",
	http.StatusNotFound:            "Not Found",
	http.StatusInternalServerError: "Internal Server Error",
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int, message string) {
	res := &ErrorResponse{StatusCode: status, Error: errors[status], Message: message}
	b, err := json.Marshal(res)
	if err != nil {
		internalError(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	fmt.Fprint(w, string(b))
}

func internalError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	fmt.Fprint(w, `{"statusCode": 500,"error": "Internal Server Error","message": "An internal server error occurred"}`)
}

func response(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(data)
}

func validate(r *http.Request, rules govalidator.MapData, data *database.User) url.Values {
	opts := govalidator.Options{
		Request: r,
		Data:    data,
		Rules:   rules,
	}

	v := govalidator.New(opts)

	return v.ValidateJSON()
}

func validationMessage(validationErrors url.Values) string {
	err := map[string]interface{}{validation: validationErrors}
	message, _ := json.Marshal(err)

	return string(message)
}
