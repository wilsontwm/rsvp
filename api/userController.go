package api

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"rsvp/models"
	"rsvp/utils"
)

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=16"`
}

type SignupInput struct {
	Name           string `json:"name" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=8,max=16"`
	RetypePassword string `json:"retype_password" validate:"required,min=8,max=16"`
	PassCode       string `json:"passcode" validate:"required"`
}

type EditProfileInput struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type EditPasswordInput struct {
	Password       string `json:"password" validate:"required,min=8,max=16"`
	RetypePassword string `json:"retype_password" validate:"required,min=8,max=16"`
}

// User login
var Login = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	input := LoginInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		resp = utils.Message(false, http.StatusInternalServerError, "Error decoding request body: "+err.Error())
		utils.Respond(w, resp)
		return
	}

	// Validate the input
	validate = validator.New()
	err = validate.Struct(input)
	if err != nil {

		resp = utils.Message(false, http.StatusUnprocessableEntity, "Validation error: "+utils.GetErrorMessages(err))

		utils.Respond(w, resp)
		return
	}

	// Login in the user
	user := &models.User{}
	resp = user.Login(input.Email, input.Password)
	utils.Respond(w, resp)
}

// User signup
var Signup = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	input := SignupInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		resp = utils.Message(false, http.StatusInternalServerError, "Error decoding request body: "+err.Error())
		utils.Respond(w, resp)
		return
	}

	// Validate the input
	validate = validator.New()
	err = validate.Struct(input)
	if err != nil {
		resp = utils.Message(false, http.StatusUnprocessableEntity, "Validation error: "+utils.GetErrorMessages(err))
		utils.Respond(w, resp)
		return
	}

	// Also do additional checking
	if input.Password != input.RetypePassword {
		resp = utils.Message(false, http.StatusUnprocessableEntity, "Validation error: Retype password does not match.")
		utils.Respond(w, resp)
		return
	}

	if input.PassCode != SignupPassCode {
		resp = utils.Message(false, http.StatusUnprocessableEntity, "Validation error: Invalid passcode.")
		utils.Respond(w, resp)
		return
	}

	// Save the data into database
	user := &models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		PassCode: input.PassCode,
	}

	// Create the account
	resp = user.Create()

	utils.Respond(w, resp)
}

// User get profile
var GetProfile = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	userId := r.Context().Value("user").(uuid.UUID)

	user := models.GetUser(userId)

	if user == nil {
		resp = utils.Message(false, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.")
		utils.Respond(w, resp)
		return
	}

	resp = utils.Message(true, http.StatusOK, "Profile retrieved")
	resp["data"] = user

	utils.Respond(w, resp)
}

// User edit profile
var EditProfile = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	userId := r.Context().Value("user").(uuid.UUID)

	user := models.GetUser(userId)

	if user == nil {
		resp := utils.Message(false, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.")
		utils.Respond(w, resp)
		return
	}

	input := EditProfileInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		resp = utils.Message(false, http.StatusInternalServerError, "Error decoding request body: "+err.Error())
		utils.Respond(w, resp)
		return
	}

	// Validate the input
	validate = validator.New()
	err = validate.Struct(input)
	if err != nil {
		resp = utils.Message(false, http.StatusUnprocessableEntity, "Validation error: "+utils.GetErrorMessages(err))
		utils.Respond(w, resp)
		return
	}

	// Save the data into database
	user.Name = input.Name
	user.Email = input.Email

	resp = user.Edit()

	utils.Respond(w, resp)
}

// User edit password
var EditPassword = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	userId := r.Context().Value("user").(uuid.UUID)

	user := models.GetUser(userId)

	if user == nil {
		resp := utils.Message(false, http.StatusUnprocessableEntity, "Something wrong has occured. Please try again.")
		utils.Respond(w, resp)
		return
	}

	input := EditPasswordInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		resp = utils.Message(false, http.StatusInternalServerError, "Error decoding request body: "+err.Error())
		utils.Respond(w, resp)
		return
	}

	// Validate the input
	validate = validator.New()
	err = validate.Struct(input)
	if err != nil {
		resp = utils.Message(false, http.StatusUnprocessableEntity, "Validation error: "+utils.GetErrorMessages(err))
		utils.Respond(w, resp)
		return
	}

	if input.Password != input.RetypePassword {
		resp = utils.Message(false, http.StatusUnprocessableEntity, "Validation error: Retype password does not match.")
		utils.Respond(w, resp)
		return
	}

	// Save the data into database
	user.Password = input.Password

	resp = user.EditPassword()

	utils.Respond(w, resp)
}
