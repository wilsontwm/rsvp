package api

import (
	"encoding/json"
	//"github.com/satori/go.uuid"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"rsvp/models"
	"rsvp/utils"
)

type AttendeesInput struct {
	Name   string   `json:"name" validate:"required"`
	Email  string   `json:"email" validate:"required,email"`
	Phone  string   `json:"phone"`
	Names  []string `json:"names"`
	Emails []string `json:"emails"`
	Phones []string `json:"phones"`
}

// Get all the attendees
var Index = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	// userId := r.Context().Value("user").(uuid.UUID)

	// user := models.GetUser(userId)

	// Authorization
	// if ok := policy.ShowInvitationFromCompany(userId, invitationId); !ok {
	// 	resp := util.Message(false, http.StatusForbidden, "You are not authorized to perform the action.", errors)
	// 	util.Respond(w, resp)
	// 	return
	// }

	attendee := &models.Attendee{}
	resp = attendee.List()

	utils.Respond(w, resp)
}

var CreateAttendees = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	input := AttendeesInput{}
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
		resp := utils.Message(false, http.StatusUnprocessableEntity, "Validation error: "+utils.GetErrorMessages(err))
		utils.Respond(w, resp)
		return
	}

	// Create the attendees
	attendee := &models.Attendee{}
	resp = attendee.CreateMultiple(input.Name, input.Email, input.Phone, input.Names, input.Emails, input.Phones)

	utils.Respond(w, resp)
}
