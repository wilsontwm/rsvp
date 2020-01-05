package utils

import (
	"encoding/json"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"reflect"
)

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
