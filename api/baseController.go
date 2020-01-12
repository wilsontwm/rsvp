package api

import (
	"github.com/joho/godotenv"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"os"
)

var validate *validator.Validate
var SignupPassCode string

func init() {
	err := godotenv.Load() //Load .env file
	if err != nil {
		log.Print("Error loading .env file", err)
	}

	SignupPassCode = os.Getenv("signup_passcode")
}
