package utils

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"os"
	"reflect"
)

// Get a session
func GetSession(store *sessions.CookieStore, w http.ResponseWriter, r *http.Request) (session *sessions.Session, err error) {
	err = godotenv.Load() //Load .env file
	sessionName := os.Getenv("session_name")
	session, err = store.Get(r, sessionName)
	return
}

// Build json message
func Message(success bool, status int, message string) map[string]interface{} {
	return map[string]interface{}{
		"success": success,
		"status":  status,
		"message": message,
	}
}

// Return json response
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	if _, ok := data["status"]; ok {
		w.WriteHeader(data["status"].(int))
	}
	json.NewEncoder(w).Encode(data)
}

// Send request to API
func SendRequest(url string, data map[string]interface{}, requestType string) (response *http.Response, err error) {
	requestBody, err := json.Marshal(data)

	request, _ := http.NewRequest(requestType, url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err = client.Do(request)

	return
}

// Build the error message
func GetErrorMessages(err error) string {
	var msg string
	for _, errz := range err.(validator.ValidationErrors) {
		// Build the custom errors here
		switch tag := errz.ActualTag(); tag {
		case "required":
			msg = errz.StructField() + " is required."
		case "email":
			msg = errz.StructField() + " is an invalid email address."
		case "min":
			if errz.Type().Kind() == reflect.String {
				msg = errz.StructField() + " must be more than or equal to " + errz.Param() + " character(s)."
			} else {
				msg = errz.StructField() + " must be larger than " + errz.Param() + "."
			}
		case "max":
			if errz.Type().Kind() == reflect.String {
				msg = errz.StructField() + " must be lesser than or equal to " + errz.Param() + " character(s)."
			} else {
				msg = errz.StructField() + " must be smaller than " + errz.Param() + "."
			}
		default:
			msg = errz.StructField() + " is invalid."
		}
	}

	return msg
}

// Initialize a page
func InitializePage(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore, data map[string]interface{}) (output map[string]interface{}, err error) {
	session, err := GetSession(store, w, r)
	errorMessages := session.Flashes("errors")
	successMessage := session.Flashes("success")
	session.Save(r, w)

	flash := map[string]interface{}{
		"errors":  errorMessages,
		"success": successMessage,
	}
	output = MergeMapString(data, flash)

	return
}

// Merge two map string interface
func MergeMapString(mp1 map[string]interface{}, mp2 map[string]interface{}) (result map[string]interface{}) {
	result = make(map[string]interface{})
	for k, v := range mp1 {
		if _, ok := mp1[k]; ok {
			result[k] = v
		}
	}

	for k, v := range mp2 {
		if _, ok := mp2[k]; ok {
			result[k] = v
		}
	}

	return result
}
